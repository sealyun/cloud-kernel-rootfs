/*
Copyright 2021 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package multiplatform //nolint:gofmt

import "fmt"

func NewOSSUtil(platform Platform) Vars {
	v := &ossutil{}
	v.setWgetURL(platform)
	return v
}

type ossutil struct {
	wgetURL string
}

func (c *ossutil) setWgetURL(platform Platform) {
	//http://gosspublic.alicdn.com/ossutil/1.7.3/ossutil64
	//https://gosspublic.alicdn.com/ossutil/1.7.5/ossutilarm64
	splatform := ""
	switch platform {
	case LinuxAmd64: //nolint:typecheck
		splatform = "64"
	case LinuxArm64: //nolint:typecheck
		splatform = "arm64"
	default:
		return
	}
	c.wgetURL = fmt.Sprintf("http://gosspublic.alicdn.com/ossutil/1.7.5/ossutil%s", splatform)
}
func (c *ossutil) FetchWgetURL() string {
	return c.wgetURL
}
func (c *ossutil) FinalShell() string {
	shell := "wget %s -O   ossutil64 &&  chmod a+x ossutil64  && mv ossutil64 /usr/bin/"
	return fmt.Sprintf(shell, c.wgetURL)
}
