#!/bin/bash 

scp -r -o StrictHostKeyChecking=no $(pwd)/src/ root@$SERVER_IP:/telegrambot-progress/

ssh -o StrictHostKeyChecking=no root@$SERVER_IP 'cd /telegrambot-progress/ && docker-compose down && docker-compose up -d --build'
