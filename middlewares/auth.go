package middlewares

import (
	"fmt"
	"github.com/alfredoptarigan/go-jwt/initializers"
	"github.com/alfredoptarigan/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	time "time"
)

func AuthProtection(c *gin.Context) {
	fmt.Println("In Middleware")

	// Get the cookies from the request
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Decode and validate it

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_TOKEN")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the expired token
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		// Find the user with token sub
		var user models.User
		initializers.DB.First(&user, "id = ?", claims["sub"])

		if user.ID == 0 {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		// Attach to request
		c.Set("user", user)

		// Continue
		fmt.Println(claims["foo"], claims["nbf"])
		c.Next()

	} else {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

}
