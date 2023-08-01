package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/watermelon03/member-api/models"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

func CheckError(c *gin.Context, err error, status int, msg string) bool {
	if err != nil {
		c.JSON(status, gin.H{
			"status":  err.Error(),
			"message": msg,
		})
		return true
	}
	return false
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(adminID, userID int) (string, error) {
	mySigningKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"adminID": adminID,
		"userID":  userID,
		"exp":     time.Now().Add(10 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}

// --------------------  API  --------------------
func TestHash() func(c *gin.Context) {
	return func(c *gin.Context) {
		hash := "$2a$14$qq4D6bUhostU9JrX5R6s.eTXbpLCzNORul5BczHoZAI9Te09iJatO"
		password := "1234"
		match := CheckPasswordHash(password, hash)

		tokenString, err := GenerateToken(1, 0)
		if CheckError(c, err, http.StatusInternalServerError, "Failed to generated token") {
			return
		}
		fmt.Println(match)
		fmt.Println(tokenString)
	}
}

func RegisterAdmin() func(c *gin.Context) {
	return func(c *gin.Context) {
		var registerBody models.RegisterBody
		err := c.ShouldBindJSON(&registerBody)
		if CheckError(c, err, http.StatusBadRequest, "Invalid request register") {
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		var adminID int
		row := DB.QueryRowContext(ctx, `SELECT adminID FROM adminaccount WHERE adminName = ?`, registerBody.Username)
		err = row.Scan(&adminID)

		if err == sql.ErrNoRows {
			hashPassword, _ := HashPassword(registerBody.Password)
			result, errInsert := DB.ExecContext(ctx, `INSERT INTO adminaccount (adminname, password, roleid, adminstate)
			VALUES (?, ?, ?, ?)`, registerBody.Username, hashPassword, registerBody.RoleID, "Available")
			if CheckError(c, errInsert, http.StatusInternalServerError, "Failed to insert adminAccount") {
				return
			}
			insertIDacc, errInsert := result.LastInsertId()
			if CheckError(c, errInsert, http.StatusInternalServerError, "Failed to retrieve insertID for adminAccount") {
				return
			}

			result, errInsert = DB.ExecContext(ctx, `INSERT INTO admininfo (adminid, firstname, lastname, 
				sex, birthday, telephone, email) VALUES (?, ?, ?, ?, ?, ?, ?)`,
				insertIDacc, registerBody.Firstname, registerBody.Lastname, registerBody.Sex,
				registerBody.Birthday, registerBody.Telephone, registerBody.Email)

			if CheckError(c, errInsert, http.StatusInternalServerError, "Failed to insert adminInfo") {
				return
			}
			insertIDinf, errInsert := result.LastInsertId()
			if CheckError(c, errInsert, http.StatusInternalServerError, "Failed to retrieve insertID for adminInfo") {
				return
			}

			c.JSON(http.StatusCreated, gin.H{
				"adminID": insertIDacc,
				"infoID":  insertIDinf,
				"status":  "ok",
				"message": "Succussfully admin register",
			})
		} else if CheckError(c, err, http.StatusInternalServerError, "Failed to check adminname exists") {
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  "error",
				"message": "AdminName already exists",
				"adminID": adminID,
			})
		}
	}
}

func LoginAdmin() func(c *gin.Context) {
	return func(c *gin.Context) {
		var loginBody models.LoginBody
		err := c.ShouldBindJSON(&loginBody)
		if CheckError(c, err, http.StatusBadRequest, "Invalid request login") {
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		row := DB.QueryRowContext(ctx, `SELECT adminID, password FROM adminaccount WHERE adminName = ?`, loginBody.Username)

		var adminID int
		var hashPassword string
		err = row.Scan(&adminID, &hashPassword)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{
				"status":  "error",
				"message": "AdminName incorrect",
			})
			return
		} else if CheckError(c, err, http.StatusInternalServerError, "Failed to check adminName exists") {
			return
		} else {
			if CheckPasswordHash(loginBody.Password, hashPassword) {
				tokenString, err := GenerateToken(adminID, 0)
				if CheckError(c, err, http.StatusInternalServerError, "Failed to generated token") {
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"adminID": adminID,
					"status":  "ok",
					"message": "Succussfully admin login",
					"token":   tokenString,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  "error",
					"message": "Login failed",
				})
			}
		}
	}
}

func RegisterUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var registerBody models.RegisterBody
		err := c.ShouldBindJSON(&registerBody)
		if CheckError(c, err, http.StatusBadRequest, "Invalid request register") {
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		var userID int
		row := DB.QueryRowContext(ctx, `SELECT userID FROM useraccount WHERE userName = ?`, registerBody.Username)
		err = row.Scan(&userID)

		if err == sql.ErrNoRows {
			hashPassword, _ := HashPassword(registerBody.Password)
			result, errInsert := DB.ExecContext(ctx, `INSERT INTO useraccount (username, password, levelID, userstate, userpoint, userqr)
			VALUES (?, ?, ?, ?, ?, ?)`, registerBody.Username, hashPassword, 4, "Available", 0, "QR CODE")
			if CheckError(c, errInsert, http.StatusInternalServerError, "Failed to insert userAccount") {
				return
			}
			insertIDacc, errInsert := result.LastInsertId()
			if CheckError(c, errInsert, http.StatusInternalServerError, "Failed to retrieve insertID for userAccount") {
				return
			}

			result, errInsert = DB.ExecContext(ctx, `INSERT INTO userinfo (userid, firstname, lastname, 
				sex, birthday, telephone, email, imagename)	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				insertIDacc, registerBody.Firstname, registerBody.Lastname, registerBody.Sex, registerBody.Birthday, registerBody.Telephone, registerBody.Email, "Null")

			if CheckError(c, errInsert, http.StatusInternalServerError, "Failed to insert userInfo") {
				return
			}
			insertIDinf, errInsert := result.LastInsertId()
			if CheckError(c, errInsert, http.StatusInternalServerError, "Failed to retrieve insertID for userInfo") {
				return
			}

			c.JSON(http.StatusCreated, gin.H{
				"adminID": insertIDacc,
				"infoID":  insertIDinf,
				"status":  "ok",
				"message": "Succussfully user register",
			})
		} else if CheckError(c, err, http.StatusInternalServerError, "Failed to check userName exists") {
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "UserName already exists",
				"userID":  userID,
			})
		}
	}
}

func LoginUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var loginBody models.LoginBody
		err := c.ShouldBindJSON(&loginBody)
		if CheckError(c, err, http.StatusBadRequest, "Invalid request login") {
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		row := DB.QueryRowContext(ctx, `SELECT userID, password FROM useraccount WHERE userName = ?`, loginBody.Username)

		var userID int
		var hashPassword string
		err = row.Scan(&userID, &hashPassword)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "UserName incorrect",
			})
			return
		} else if CheckError(c, err, http.StatusInternalServerError, "Failed to check userName exists") {
			return
		} else {
			if CheckPasswordHash(loginBody.Password, hashPassword) {
				tokenString, err := GenerateToken(0, userID)
				if CheckError(c, err, http.StatusInternalServerError, "Failed to generated token") {
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"userID":  userID,
					"status":  "ok",
					"message": "Succussfully user login",
					"token":   tokenString,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  "error",
					"message": "Login failed",
				})
			}
		}
	}
}

// ------------------------------------------
type StructB struct {
	FieldA string `form:"field_a"`
	FieldB string `form:"field_b"`
	FieldC string `json:"field_C"`
}

func GetDataB() func(c *gin.Context) {
	return func(c *gin.Context) {
		var b StructB
		c.Bind(&b)
		c.JSON(200, gin.H{
			"a": b.FieldA,
			"b": b.FieldB,
		})
	}

}

func UploadFormData() func(c *gin.Context) {
	return func(c *gin.Context) {
		form, _ := c.MultipartForm()

		// Retrieve the files
		files := form.File["files"]

		// Process each file
		for _, file := range files {
			// Save the file to disk or perform any other desired operation
			if err := c.SaveUploadedFile(file, "images/"+file.Filename); err != nil {
				c.String(500, "Failed to upload file")
				return
			}
		}
		var b StructB
		c.ShouldBindJSON(&b)

		c.JSON(200, gin.H{
			"status": "Files uploaded successfully",
			"a":      b.FieldA,
			"b":      b.FieldB,
			"c":      b.FieldC,
		})
	}
}
