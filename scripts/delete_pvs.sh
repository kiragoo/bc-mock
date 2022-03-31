#!/bin/bash
kubectl get pv | awk 'NR!=1 {print $1}' > .tmp.pv

for line in `cat .tmp.pv`
    do kubectl delete pv "$line"
done
rm -rf .tmp.pv