#! /bin/bash
API_TOKEN="xxx,xxxxxx"
DOMAIN="xx"
DOMAIN_ID="xxxxx"
RECORD_ID="xxxxx"
current_ip=$(curl -sSL ipinfo.io/ip)
echo $current_ip
response=$(curl -s -X POST https://dnsapi.cn/Record.Ddns -d "login_token=${API_TOKEN}&format=json&sub_domain=${DOMAIN}&domain_id=${DOMAIN_ID}&record_line=默认&record_id=${RECORD_ID}&value=${current_ip}")
echo $response
