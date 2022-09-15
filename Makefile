build:
	docker build -t grafana-alert .

up:
	docker run --rm --restart=unless-stopped -e BOT_TOKEN=5103252943:AAG3Xd5yGwFtXTnZ2HnWrcyLlOatp4k2Ook -e CHAT_ID=-1001750896977 -p 1323:1323 --name=alert -d grafana-alert

down:
	docker stop alert

rm:
	docker rmi grafana-alert && docker system prune

logs:
	docker logs -f alert
