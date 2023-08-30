package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opshive/go-td/controllers/helpers"
	"github.com/opshive/go-td/models"
	"github.com/opshive/go-td/pkg/clients"
	"gopkg.in/yaml.v2"
)

// @Summary Deploy app with helm
// @Description Deploy a new app by providing a JSON object with the "name" property.
// @Tags helm chart
// @Accept json
// @Produce json
// @Param name body string true "app Name"
// @Success 200 {string} string "app created successfully"
// @Router /appCreate [post]
func HelmlDeploy(c *gin.Context) {
	var userInput struct {
		ImageRepo string `json:"imageRepo"`
		SubDomain string `json:"subDomain"`
		Name      string `json:"name"`
	}
	// implement redis cache

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	nginxConfig := models.NginxConfig{
		ReplicaCount: 1,
		Image: models.ImageConfig{
			Repository: userInput.ImageRepo,
			PullPolicy: "IfNotPresent",
			Tag:        "",
		},
		ImagePullSecrets: []interface{}{},
		NameOverride:     "",
		FullnameOverride: "",
		ServiceAccount: models.ServiceAccountConfig{
			Create:      true,
			Annotations: map[string]string{},
			Name:        "",
		},
		PodAnnotations:     map[string]string{},
		PodSecurityContext: map[string]interface{}{},
		SecurityContext:    map[string]interface{}{},
		Service: struct {
			Type string `yaml:"type"`
			Port int    `yaml:"port"`
		}{
			Type: "ClusterIP",
			Port: 80,
		},
		Ingress: models.IngressConfig{
			Enabled:     false,
			ClassName:   "",
			Annotations: map[string]string{},
			Hosts: []models.HostConfig{
				{
					Host: fmt.Sprintf("https://%s.opshive.io", userInput.SubDomain),
					Paths: []models.PathConfig{
						{
							Path:     "/",
							PathType: "ImplementationSpecific",
						},
					},
				},
			},
			TLS: []interface{}{},
		},
		Resources: struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		}{
			Limits: map[string]string{
				"cpu":    "100m",
				"memory": "128Mi",
			},
			Requests: map[string]string{
				"cpu":    "100m",
				"memory": "128Mi",
			},
		},
		Autoscaling: struct {
			Enabled                        bool `yaml:"enabled"`
			MinReplicas                    int  `yaml:"minReplicas"`
			MaxReplicas                    int  `yaml:"maxReplicas"`
			TargetCPUUtilizationPercentage int  `yaml:"targetCPUUtilizationPercentage"`
		}{
			Enabled:                        false,
			MinReplicas:                    1,
			MaxReplicas:                    100,
			TargetCPUUtilizationPercentage: 80,
		},
		NodeSelector: map[string]string{},
		Tolerations:  []interface{}{},
		Affinity:     map[string]interface{}{},
	}

	// Marshal the struct into YAML
	yamlData, err := yaml.Marshal(&nginxConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// You can then write the YAML data to a file or print it to the console
	// For example, to write to a file:
	tmpfile, err := os.CreateTemp("./", "nginx-config-*.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file

	if _, err = tmpfile.Write(yamlData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write to temporary file"})
		return
	}
	if err := tmpfile.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close temporary file"})
		return
	}
	kubeconfigPath, err := clients.GetKubeconfigPath(os.Getenv("APP_ENV"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Kubernetes cluster", "details": err.Error()})
		return
	}

	helmChartPath := "./internals/charts/nginx"
	output, err := helpers.RetryHelmCommand("helm", []string{"install", userInput.Name, helmChartPath, "--values", tmpfile.Name()}, kubeconfigPath, 3, 5*time.Second)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deploy helm chart after multiple attempts", "details": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "app created successfully", "details": string(output)})
	}
}

// @Summary Deploy app with helm
// @Description app status "name" property.
// @Tags helm chart
// @Produce json
// @Success 200 {string} string "app status"
// @Router /appGet [get]
func GetChart(c *gin.Context) {
	kubeconfigPath, err := clients.GetKubeconfigPath(os.Getenv("APP_ENV"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Kubernetes cluster", "details": err.Error()})
		return
	}
	// cmd := exec.Command("helm", "list", "-o", "json")
	// cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)

	output, err := helpers.RetryHelmCommand("helm", []string{"list", "-o", "json"}, kubeconfigPath, 3, 5*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list helm chart", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sucess", "details": output})

}

// @Summary Deploy app with helm
// @Description Delete deployed app by providing a JSON object with the "name" property.
// @Tags helm chart
// @Accept json
// @Produce json
// @Param name body string true "app Name"
// @Success 200 {string} string "app Deleted successfully"
// @Router /appDelete [post]
func ChartDelete(c *gin.Context) {
	var userInput struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kubeconfigPath, err := clients.GetKubeconfigPath(os.Getenv("APP_ENV"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Kubernetes cluster", "details": err.Error()})
		return
	}
	// cmd := exec.Command("helm", "uninstall", userInput.Name)
	// cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)

	output, err := helpers.RetryHelmCommand("helm", []string{"uninstall", userInput.Name}, kubeconfigPath, 3, 5*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete helm chart", "details": string(output)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "app deleted successfully"})
}
