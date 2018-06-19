dnatip="49.4.66.120"
dnatid="67d11284-9ddc-4edd-aa3e-ce544970e89f"
elbip="49.4.50.23"
elbid="57e9121c-3df1-4b20-8662-7e7c1a480211"
elbport="8088"
sfsname="cce-sfs-jiehlm60-ghuf"
rm -rf tmp/
mkdir tmp
cp ingress.yaml tmp/ingress.yaml
sed -i "s|{{elbip}}|$elbip|g" tmp/ingress.yaml
sed -i "s|{{elbid}}|$elbid|g" tmp/ingress.yaml
sed -i "s|{{elbport}}|$elbport|g" tmp/ingress.yaml
for (( i=1; i<=${1}; i++ ))
do
    echo "      - backend:" >> tmp/ingress.yaml
    echo "          serviceName: demotest${i}-nodeport" >> tmp/ingress.yaml
    echo "          servicePort: 8088" >>  tmp/ingress.yaml
    echo "        path: \"/api/v1/demotest${i}\"" >>  tmp/ingress.yaml
    echo "        property:">>  tmp/ingress.yaml
    echo "          ingress.beta.kubernetes.io/url-match-mode: STARTS_WITH" >> tmp/ingress.yaml
    ((port=33100+$i))
    ((port2=$port+1000))
    cp demotest.yaml tmp/demotest${i}.yaml
    sed -i "s|{{dnatip}}|$dnatip|g" tmp/demotest${i}.yaml
    sed -i "s|{{dnatid}}|$dnatid|g" tmp/demotest${i}.yaml
    sed -i "s|{{elbip}}|$elbip|g" tmp/demotest${i}.yaml
    sed -i "s|{{elbid}}|$elbid|g" tmp/demotest${i}.yaml
    sed -i "s|{{sfsname}}|$sfsname|g" tmp/demotest${i}.yaml
    sed -i "s|{{servername}}|demotest${i}|g" tmp/demotest${i}.yaml
    sed -i "s|{{port}}|$port|g" tmp/demotest${i}.yaml
    sed -i "s|{{port2}}|$port2|g" tmp/demotest${i}.yaml
    echo "kubectl delete -f demotest${i}.yaml" >> tmp/delete.sh
    echo "kubectl create -f demotest${i}.yaml" >> tmp/create.sh
done

