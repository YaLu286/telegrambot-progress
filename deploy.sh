#!/bin/bash 

scp -r -o StrictHostKeyChecking=no $(pwd)/src/ clamsmil@$SERVER_IP:/home/clamsmil/telegrambot-progress/

ssh -o StrictHostKeyChecking=no clamsmil@$SERVER_IP 'cd /home/clamsmil/telegrambot-progress/ && docker-compose down && docker-compose up -d --build'
