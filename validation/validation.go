/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package validation

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"intel/tac/v1/constants"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	stringReg = regexp.MustCompile("(^[a-zA-Z0-9_ \\/.-]*$)")
	emailReg  = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+\/=?^_'{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)
	// Regex to validate TA API key. Key should contain characters between a-z, A-Z, 0-9
	// and should be of size between 30 and 128
	apiKeyRegex           = regexp.MustCompile(`^[a-zA-Z0-9]{30,128}$`)
	subscriptionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-\_]{1,62}[a-zA-Z0-9]$`)
	tagReg                = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-\_]{1,62}[a-zA-Z0-9]$`)
	tagValueReg           = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-\_]{1,62}[a-zA-Z0-9]$`)
	policyNameRegex       = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{1,62}[a-zA-Z0-9]$`)
	requestIdRegex        = regexp.MustCompile(`^[a-zA-Z0-9_ \/.-]{1,128}$`)
)

// ValidateStrings method is used to validate input strings
func ValidateStrings(strings []string) error {
	for _, stringValue := range strings {
		if !stringReg.MatchString(stringValue) {
			return errors.New("Invalid string formatted input")
		}
	}
	return nil
}

func ValidateEmailAddress(email string) error {
	if !emailReg.Match([]byte(email)) {
		logrus.Error("Invalid email id provided")
		return errors.New("Invalid email id provided")
	}

	return nil
}

func ValidatePath(path string) (string, error) {
	c := filepath.Clean(path)
	r, err := filepath.EvalSymlinks(c)
	if err != nil {
		return c, fmt.Errorf("%s: %s", constants.ErrorInvalidPath, path)
	}
	return r, nil
}

func ValidateSize(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.Size() > constants.MaxPolicyFileSize {
		return fmt.Errorf("%s: %d", constants.ErrorInvalidSize, fi.Size())
	}
	return nil
}

func ValidateTrustAuthorityAPIKey(apiKey string) error {
	if strings.TrimSpace(apiKey) == "" {
		return errors.Errorf("%s config variable needs to be set with a proper API Key before using CLI", constants.TrustAuthApiKeyEnvVar)
	}
	if !apiKeyRegex.MatchString(apiKey) {
		return errors.New("Invalid API key found in configuration file. Please update it with a valid API key.")
	}
	return nil
}

func ValidateApiClientName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("ApiClient name cannot be empty")
	}
	if !subscriptionNameRegex.Match([]byte(name)) {
		return errors.New("ApiClient name should be alphanumeric and start with an alphanumeric character with " +
			"_ or - as separator and should be at most 64 characters long")
	}
	return nil
}

func ValidateTagName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("Tag name cannot be empty")
	}
	if !tagReg.Match([]byte(name)) {
		return errors.New("Tag name should be alphanumeric and start with an alphanumeric character with " +
			"_ or - as separator and should be at most 64 characters long")
	}
	return nil
}

func ValidateTagValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("Tag value cannot be empty")
	}
	if !tagValueReg.Match([]byte(value)) {
		return errors.New("Tag value should be alphanumeric and start with an alphanumeric character with " +
			"_ or - as separator and should be at most 64 characters long")
	}
	return nil
}

func ValidatePolicyName(policyName string) error {
	if strings.TrimSpace(policyName) == "" {
		return errors.New("Policy name cannot be empty")
	}
	if !policyNameRegex.Match([]byte(policyName)) {
		return errors.New("Policy name is invalid. Policy name should be alpha numeric and have minimum 3 characters with no spaces between words (" +
			"use \"_\" or \"-\" as separators) and should not be more than 64 characters")
	}
	return nil
}

func ValidateRequestId(requestId string) error {
	if strings.TrimSpace(requestId) != "" && !requestIdRegex.Match([]byte(requestId)) {
		return errors.New("Request ID should be at most 128 characters long and should contain only " +
			"alphanumeric characters, _, space, - or \\")
	}
	return nil
}
