package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type ProjectAutoApiService struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	AppConfig struct {
		Endpoint string `json:"endpoint"`
		DbSchema string `json:"db_schema"`
	} `json:"app_config"`
	App struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"app"`
	ServiceApiKeys []struct {
		Tags string `json:"tags"`
		Name string `json:"name"`
	} `json:"service_api_keys"`
	Protocol string `json:"protocol"`
	Endpoint string `json:"endpoint"`
	RestUrl  string `json:"restUrl"`
	Project  struct {
		Ref string `json:"ref"`
	} `json:"project"`
	DefaultApiKey string `json:"defaultApiKey"`
	ServiceApiKey string `json:"serviceApiKey"`
}

func (a *Api) getProjectApi(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	projectRef := c.Param("ref")
	proj, err := a.queries.GetProjectByRef(c, projectRef)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// Get real keys - use stored keys if available, otherwise generate from jwt_secret
	var anonKey, serviceKey string

	if proj.AnonKey.Valid && proj.AnonKey.String != "" {
		anonKey = proj.AnonKey.String
	} else {
		// Generate anon key from JWT secret
		var genErr error
		anonKey, genErr = a.generateProjectJWT(proj.JwtSecret, "anon")
		if genErr != nil {
			c.JSON(500, gin.H{"error": "Failed to generate anon key"})
			return
		}
	}

	if proj.ServiceRoleKey.Valid && proj.ServiceRoleKey.String != "" {
		serviceKey = proj.ServiceRoleKey.String
	} else {
		// Generate service_role key from JWT secret
		var genErr error
		serviceKey, genErr = a.generateProjectJWT(proj.JwtSecret, "service_role")
		if genErr != nil {
			c.JSON(500, gin.H{"error": "Failed to generate service key"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"autoApiService": ProjectAutoApiService{
			Id:   0,
			Name: "Default API",
			AppConfig: struct {
				Endpoint string `json:"endpoint"`
				DbSchema string `json:"db_schema"`
			}{
				Endpoint: fmt.Sprintf("%s.supamanager.io", proj.ProjectRef),
				DbSchema: "public",
			},
			App: struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			}{
				Id:   1,
				Name: "Auto API",
			},
			ServiceApiKeys: []struct {
				Tags string `json:"tags"`
				Name string `json:"name"`
			}{
				{
					Tags: "anon",
					Name: "anon key",
				},
				{
					Tags: "service_role",
					Name: "service_role key",
				},
			},
			Protocol: "https",
			Endpoint: fmt.Sprintf("%s.supamanager.io", proj.ProjectRef),
			RestUrl:  fmt.Sprintf("https://%s.supamanager.io/rest/v1/", proj.ProjectRef),
			Project: struct {
				Ref string `json:"ref"`
			}{
				Ref: proj.ProjectRef,
			},
			DefaultApiKey: anonKey,
			ServiceApiKey: serviceKey,
		},
	})
}

// generateProjectJWT generates a Supabase-compatible JWT token
// role can be "anon" or "service_role"
func (a *Api) generateProjectJWT(jwtSecret string, role string) (string, error) {
	// Supabase JWT claims
	claims := jwt.MapClaims{
		"iss":  "supamanager",
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().AddDate(10, 0, 0).Unix(), // 10 years expiry (like Supabase)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
