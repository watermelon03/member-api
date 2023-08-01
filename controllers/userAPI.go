package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/watermelon03/member-api/models"
)

func GetUserProfile() func(c *gin.Context) {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(float64)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		row := DB.QueryRowContext(ctx, `SELECT * FROM useraccount INNER JOIN userinfo 
		ON useraccount.userID = userinfo.userID WHERE useraccount.userID = ?`, userID)

		var userAccount models.UserAccount
		var userInfo models.UserInfo
		var blank string
		var levelID int

		err := row.Scan(&userAccount.UserID,
			&userAccount.UserName,
			&blank,
			&levelID,
			&userAccount.UserState,
			&userAccount.UserPoint,
			&blank,
			&userAccount.AccountUpdateDate,
			&blank,
			&userInfo.UserID,
			&userInfo.Firstname,
			&userInfo.Lastname,
			&userInfo.Sex,
			&userInfo.Birthday,
			&userInfo.Telephone,
			&userInfo.Email,
			&userInfo.RegisterDate,
			&userInfo.InfoUpdateDate,
			&userInfo.ImageName)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		} else if CheckError(c, err, http.StatusInternalServerError, "Failed to get user") {
			return
		}
		row = DB.QueryRowContext(ctx, `SELECT levelName, levelImage FROM userlevel WHERE levelID = ?`, levelID)
		row.Scan(&userAccount.LevelName, &userAccount.LevelImage)

		fakePassword := "###*****###*****###"
		userAccount.Password = fakePassword

		c.JSON(http.StatusOK, gin.H{
			"userAccount": userAccount,
			"userInfo":    userInfo,
		})
	}
}

func UpdateUserImage() func(c *gin.Context) {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(float64)
		userIDStr := strconv.Itoa(int(userID))

		form, _ := c.MultipartForm()
		files := form.File["image-files"]

		fileName := userIDStr + "_" + files[0].Filename
		filePath := "user-images/" + fileName

		err := c.SaveUploadedFile(files[0], filePath)
		if CheckError(c, err, http.StatusInternalServerError, "Failed to save file image") {
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
		defer cancel()

		result, err := DB.ExecContext(ctx, `UPDATE userinfo SET imagename = ? WHERE userID = ?`, fileName, userID)
		if CheckError(c, err, http.StatusInternalServerError, "Failed to update user image") {
			return
		}
		rowAff, err := result.RowsAffected()
		if CheckError(c, err, http.StatusInternalServerError, "Failed to returns the number of rows affected") {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"rowAffected": rowAff,
			"message":     "Succussfully update user image",
		})
	}
}

func UpdateUserPassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(float64)

		var userAccount models.UserAccount
		err := c.ShouldBindJSON(&userAccount)
		if CheckError(c, err, http.StatusBadRequest, "Invalid request update user password") {
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
		defer cancel()

		var hashPassword string
		row := DB.QueryRowContext(ctx, `SELECT password FROM useraccount WHERE userID = ?`, userID)
		row.Scan(&hashPassword)

		if CheckPasswordHash(userAccount.Password, hashPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password same as old one",
			})
			return
		} else {
			newHashPassword, _ := HashPassword(userAccount.Password)
			result, err := DB.ExecContext(ctx, `UPDATE useraccount SET password = ? WHERE userID = ?`, newHashPassword, userID)
			if CheckError(c, err, http.StatusInternalServerError, "Failed to update user password") {
				return
			}
			rowAff, err := result.RowsAffected()
			if CheckError(c, err, http.StatusInternalServerError, "Failed to returns the number of rows affected") {
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"rowAffected": rowAff,
				"message":     "Succussfully update user password",
			})
		}
	}
}

func UpdateUserInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(float64)

		var userInfo models.UserInfo
		err := c.ShouldBindJSON(&userInfo)
		if CheckError(c, err, http.StatusBadRequest, "Invalid request update user infomation") {
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
		defer cancel()

		result, err := DB.ExecContext(ctx, `UPDATE userinfo SET firstname = ?, lastname = ?, 
		sex = ?, birthday = ?, telephone = ?, email = ? WHERE userID = ?`,
			userInfo.Firstname, userInfo.Lastname, userInfo.Sex, userInfo.Birthday, userInfo.Telephone, userInfo.Email, userID)

		if CheckError(c, err, http.StatusInternalServerError, "Failed to update user infomation") {
			return
		}
		rowAff, err := result.RowsAffected()
		if CheckError(c, err, http.StatusInternalServerError, "Failed to returns the number of rows affected") {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"rowAffected": rowAff,
			"message":     "Succussfully update user infomation",
		})
	}
}
