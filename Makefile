init:
	docker-compose up -d

run-local:
	docker-compose exec env /bin/bash -c 'go run ./cmd/*.go run-local'

deploy:
	gcloud functions deploy BlogUpdateArticle --project suzuito-$${ENV} --region asia-northeast1 --max-instances 1 --timeout 60s --runtime go113 --env-vars-file=env-$${ENV}.yml --trigger-resource suzuito-$${ENV}-blog1-article --trigger-event google.storage.object.finalize