package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/watermelon03/member-api/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetAdminAll() func(c *gin.Context) {
	return func(c *gin.Context) {
		adminID := c.MustGet("adminID").(float64)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		results, err := DB.QueryContext(ctx, `SELECT * FROM adminaccount`)
		if CheckError(c, err, http.StatusInternalServerError, "Failed to get all admin") {
			return
		}
		defer results.Close()

		var adminAccounts []models.AdminAccount
		for results.Next() {
			var adminAccount models.AdminAccount
			var roleID int
			var blank string
			results.Scan(&adminAccount.AdminID,
				&adminAccount.AdminName,
				&blank,
				&roleID,
				&adminAccount.AdminState,
				&adminAccount.AccountUpdateDate)

			row := DB.QueryRowContext(ctx, `SELECT roleName FROM adminrole WHERE roleID = ?`, roleID)
			row.Scan(&adminAccount.RoleName)

			fakePassword := "******####******"
			adminAccount.Password = fakePassword

			adminAccounts = append(adminAccounts, adminAccount)
		}
		c.JSON(http.StatusOK, gin.H{
			"adminID":    adminID,
			"AllAccount": adminAccounts,
			"status":     "ok",
			"message":    "Succussfully get all adminAccounts",
		})
	}
}

func GetAdminProfile() func(c *gin.Context) {
	return func(c *gin.Context) {
		adminID := c.MustGet("adminID").(float64)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		row := DB.QueryRowContext(ctx, `SELECT * FROM adminaccount INNER JOIN admininfo 
		ON adminaccount.adminID = admininfo.adminID WHERE adminaccount.adminID = ?`, adminID)

		var adminAccount models.AdminAccount
		var adminInfo models.AdminInfo
		var blank string
		var roleID int

		err := row.Scan(&adminAccount.AdminID,
			&adminAccount.AdminName,
			&blank,
			&roleID,
			&adminAccount.AdminState,
			&adminAccount.AccountUpdateDate,
			&blank,
			&adminInfo.AdminID,
			&adminInfo.Firstname,
			&adminInfo.Lastname,
			&adminInfo.Sex,
			&adminInfo.Birthday,
			&adminInfo.Telephone,
			&adminInfo.Email,
			&adminInfo.RegisterDate,
			&adminInfo.InfoUpdateDate)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Admin not found",
			})
			return
		} else if CheckError(c, err, http.StatusInternalServerError, "Failed to get admin") {
			return
		}
		row = DB.QueryRowContext(ctx, `SELECT roleName FROM adminrole WHERE roleID = ?`, roleID)
		row.Scan(&adminAccount.RoleName)

		fakePassword := "******####******"
		adminAccount.Password = fakePassword

		c.JSON(http.StatusOK, gin.H{
			"adminAccount": adminAccount,
			"adminInfo":    adminInfo,
			"status":       "ok",
			"message":      "Succussfully get adminInfo",
		})
	}
}

func UpdateAdminPassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		adminID := c.MustGet("adminID").(float64)

		var adminAccount models.AdminAccount
		err := c.ShouldBindJSON(&adminAccount)
		if CheckError(c, err, http.StatusBadRequest, "Invalid request update admin password") {
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
		defer cancel()

		var hashPassword string
		row := DB.QueryRowContext(ctx, `SELECT password FROM adminaccount WHERE adminID = ?`, adminID)
		row.Scan(&hashPassword)

		if CheckPasswordHash(adminAccount.Password, hashPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Password same as old one",
			})
			return
		} else {
			newHashPassword, _ := HashPassword(adminAccount.Password)
			result, err := DB.ExecContext(ctx, `UPDATE adminaccount SET password = ? WHERE adminID = ?`, newHashPassword, adminID)
			if CheckError(c, err, http.StatusInternalServerError, "Failed to update admin password") {
				return
			}
			rowAff, err := result.RowsAffected()
			if CheckError(c, err, http.StatusInternalServerError, "Failed to returns the number of rows affected") {
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"rowAffected": rowAff,
				"status":      "ok",
				"message":     "Succussfully update admin password",
			})
		}
	}
}

func UpdateAdminInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		adminID := c.MustGet("adminID").(float64)

		var adminInfo models.AdminInfo
		err := c.ShouldBindJSON(&adminInfo)
		if CheckError(c, err, http.StatusBadRequest, "Invalid request update admin infomation") {
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
		defer cancel()

		result, err := DB.ExecContext(ctx, `UPDATE admininfo SET firstname = ?, lastname = ?, 
		sex = ?, birthday = ?, telephone = ?, email = ? WHERE adminID = ?`,
			adminInfo.Firstname, adminInfo.Lastname, adminInfo.Sex, adminInfo.Birthday, adminInfo.Telephone, adminInfo.Email, adminID)

		if CheckError(c, err, http.StatusInternalServerError, "Failed to update infomation") {
			return
		}
		rowAff, err := result.RowsAffected()
		if CheckError(c, err, http.StatusInternalServerError, "Failed to returns the number of rows affected") {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"rowAffected": rowAff,
			"status":      "ok",
			"message":     "Succussfully update admin infomation",
		})
	}
}

func GetUserAll() func(c *gin.Context) {
	return func(c *gin.Context) {
		adminID := c.MustGet("adminID").(float64)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		results, err := DB.QueryContext(ctx, `SELECT * FROM useraccount`)
		if CheckError(c, err, http.StatusInternalServerError, "Failed to get all user") {
			return
		}
		defer results.Close()

		var userAccounts []models.UserAccount
		for results.Next() {
			var userAccount models.UserAccount
			var levelID int
			var blank string
			results.Scan(&userAccount.UserID,
				&userAccount.UserName,
				&blank,
				&levelID,
				&userAccount.UserState,
				&userAccount.UserPoint,
				&userAccount.UserQR,
				&userAccount.AccountUpdateDate)

			row := DB.QueryRowContext(ctx, `SELECT levelName, levelImage FROM userLevel WHERE levelID = ?`, levelID)
			row.Scan(&userAccount.LevelName, &userAccount.LevelImage)

			fakePassword := "******####******"
			userAccount.Password = fakePassword

			userAccounts = append(userAccounts, userAccount)
		}
		c.JSON(http.StatusOK, gin.H{
			"adminID":    adminID,
			"AllAccount": userAccounts,
			"status":     "ok",
			"message":    "Succussfully get all userAccounts",
		})
	}
}
