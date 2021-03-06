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
package templates

import (
	"bytes"
	"text/template"
)

type Templates interface {
	Template() string
	TemplateConvert() string
}

func templateFromContent(templateContent string, param map[string]interface{}) (string, error) {
	tmpl, err := template.New("text").Parse(templateContent)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, param)
	bs := buffer.Bytes()
	if nil != bs && len(bs) > 0 {
		return string(bs), nil
	}
	return "", err
}
