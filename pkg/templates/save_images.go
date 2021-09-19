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

const (
	SaveImageDocker = `#!/bin/bash
echo "127.0.0.1 sea.hub" >> /etc/hosts
mv ../images/registry.tar ../images/registry_cache.tar
sh init-registry.sh 5000 /var/lib/registry
images=$(docker images | grep -v registry | grep -v REPOSITORY | grep k8s.gcr.io | awk '{print $1":"$2}')
for i in $images ; do
    echo "pushing sea.hub:5000/library/${i##k8s.gcr.io/}"
    docker tag $i sea.hub:5000/library/${i##k8s.gcr.io/}
    docker push sea.hub:5000/library/${i##k8s.gcr.io/}
done

images=$(docker images | grep -v registry | grep -v REPOSITORY | grep -v k8s.gcr.io | grep -v library | awk '{print $1":"$2}')
for i in $images ; do
    echo "pushing sea.hub:5000/$i"
    docker tag $i sea.hub:5000/$i
    docker push sea.hub:5000/$i
done

cp -rf /var/lib/registry/docker ../registry
mv ../images/registry_cache.tar ../images/registry.tar
`
)
