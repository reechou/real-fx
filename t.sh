
# maybe more powerful
# for mac (sed for linux is different)
grep "x-http-server" * -R | grep -v Godeps | awk -F: '{print $1}' | sort | uniq | xargs sed -i '' 's#x-http-server#real-fx#g'
