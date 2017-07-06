/*
Copyright (c) 2017, UPMC Enterprises
All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name UPMC Enterprises nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL UPMC ENTERPRISES BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
*/

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

var (
	argAction       = flag.String("action", "", "action to perform (e.g. Create repository or snapshot")
	argS3BucketName = flag.String("s3-bucket-name", "", "name of s3 bucket")
	argElasticURL   = flag.String("elastic-url", "", "full dns url to elasticsearch")
	argUsername     = flag.String("auth-username", "", "Authentication username")
	argPassword     = flag.String("auth-password", "", "Authentication password")
)

// CreateSnapshotRepository creates a repository to place snapshots
func CreateSnapshotRepository() {
	logrus.Info("About to create Snapshot Repository...")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("%s/_snapshot/%s", *argElasticURL, *argS3BucketName)
	body := fmt.Sprintf("{ \"type\": \"s3\", \"settings\": { \"bucket\": \"%s\" } }", *argS3BucketName)
	req, err := http.NewRequest("PUT", url, strings.NewReader(body))

	// if authentication is specified, provide Auth to Client
	if *argUsername != "" && *argPassword != "" {
		logrus.Infof("Using basic Auth Credentials %s", *argUsername)
		req.SetBasicAuth(*argUsername, *argPassword)
	}

	resp, err := client.Do(req)

	// Some other type of error?
	if err != nil {
		logrus.Error("Error attempting to create snapshot repository: ", err)
		return
	}

	// Non 2XX status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		logrus.Errorf("Error creating snapshot repository [httpstatus: %d][url: %s][body: %s] ", resp.StatusCode, url, string(body))
		return
	}

	logrus.Infof("Created snapshot repository!")
}

// CreateSnapshot makes a snapshot of all indexes
func CreateSnapshot() {
	logrus.Info("About to create snapshot...")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("%s/_snapshot/%s/snapshot_%s?wait_for_completion=true", *argElasticURL, *argS3BucketName, fmt.Sprintf(time.Now().Format("2006-01-02-15-04-05")))

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		logrus.Error("Error attempting to create snapshot: ", err)
		return
	}

	// if authentication is specified, provide Auth to Client
	if *argUsername != "" && *argPassword != "" {
		logrus.Infof("Using basic Auth Credentials %s", *argUsername)
		req.SetBasicAuth(*argUsername, *argPassword)
	}

	resp, err := client.Do(req)

	// Some other type of error?
	if err != nil {
		logrus.Error("Error attempting to create snapshot: ", err)
		return
	}

	// Non 2XX status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		logrus.Errorf("Error creating snapshot [httpstatus: %d][url: %s] %s", resp.StatusCode, url, string(body))
		return
	}

	logrus.Infof("Created snapshot!")
}

func main() {
	flag.Parse()
	log.Println("[elasticsearch-cron] is up and running!", time.Now())

	// Validate
	if *argElasticURL == "" {
		logrus.Fatal("Missing ElasticURL parameter!")
	}

	if *argS3BucketName == "" {
		logrus.Fatal("Missing S3 Bucket Name parameter!")
	}

	switch *argAction {
	case "create-repository":
		CreateSnapshotRepository()
		break
	case "snapshot":
		CreateSnapshot()
		break
	default:
		logrus.Infof("Command passed [%s] not recognized.", *argAction)
	}
}
