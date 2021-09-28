package core

import "regexp"

type Image struct {

	/**
	 *	镜像名称
	 */
	SourceImage    string `json:"source_image"`
	SourceAuth     string `json:"source_auth"`
	SourceUser     string `json:"source_user"`
	SourcePassword string `json:"source_password"`

	/**
	 * minioBucketName
	 */
	TargetImage    string `json:"target_image"`
	TargetAuth     string `json:"target_auth"`
	TargetUser     string `json:"target_user"`
	TargetPassword string `json:"target_password"`
}

func (image *Image) EmptySourceAuth() bool {
	if image.SourceUser == "" && image.SourcePassword == "" && image.SourceAuth == "" {
		return true
	}
	return false
}

func (image *Image) EmptyTargetAuth() bool {
	if image.TargetUser == "" && image.TargetPassword == "" && image.TargetAuth == "" {
		return true
	}
	return false
}

func (image *Image) BuildRegistryAuth(registryAuth map[string]string) {
	re := regexp.MustCompile("[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\\.?")
	sourceRegistry := re.FindString(image.SourceImage)
	targetRegistry := re.FindString(image.TargetImage)
	if val, ok := registryAuth[sourceRegistry]; ok {
		image.SourceAuth = val
	}
	if val, ok := registryAuth[targetRegistry]; ok {
		image.TargetAuth = val
	}
}
