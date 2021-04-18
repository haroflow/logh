# Log Highlighter

This tool colorizes the output based on regular expressions, making it easier to follow logs at a glance.
![logh](https://user-images.githubusercontent.com/4776931/115137317-d8edc180-9ffb-11eb-8542-c84260cafbc0.png)

# Requirements
- go 1.16

# Installing
```
git clone https://github.com/haroflow/logh
cd logh
go build cli/logh.go
# Copy logh to a folder in your PATH
```

# Using

Pass data to `logh` via stdin.
Pass regular expressions as arguments to highlight the text. Each argument will get a different color (up to 10 at the moment).

Examples:
```
# View help
logh

# Highlights "GET" with red, "POST" with green
tail -f /var/log/httpd/access.log | logh GET POST

# Highlights the whole line containing "word"
cat textfile.txt | logh '.*word.*'

# Ignore case. Matches get, GET, Get, etc.
cat textfile.txt | logh -i get
```
