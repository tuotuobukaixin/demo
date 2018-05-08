EIP="117.78.26.113"
sfsname="cce-sfs-jgg3vxr1-xfsp"
rm -rf tmp/
mkdir tmp
for (( i=1; i<=${1}; i++ ))
do
    ((port=31100+$i))
    cp gameserver.yaml tmp/gameserver${i}.yaml
    sed -i "s|{{EIP}}|$EIP|g" tmp/gameserver${i}.yaml
    sed -i "s|{{sfsname}}|$sfsname|g" tmp/gameserver${i}.yaml
    sed -i "s|{{servername}}|gameserver${i}|g" tmp/gameserver${i}.yaml
    sed -i "s|{{port}}|$port|g" tmp/gameserver${i}.yaml
done