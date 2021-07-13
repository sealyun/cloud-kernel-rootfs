#!/bin/bash
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

