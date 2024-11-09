ENDPOINT_URL=http://localhost:4566
PROFILE=localstack
QUEUE_URL=http://sqs.ap-northeast-1.localhost.localstack.cloud:4566/000000000000/message-queue

.PHONY: sqs-list
sqs-list:
	aws --endpoint-url=$(ENDPOINT_URL) --profile $(PROFILE)\
		sqs list-queues

.PHONY: sqs-receive
sqs-receive:
	aws --endpoint-url=$(ENDPOINT_URL) --profile $(PROFILE)\
		sqs receive-message --queue-url $(QUEUE_URL) --visibility-timeout 1

.PHONY: sqs-purge
sqs-purge:
	while true; do \
	  MESSAGES=$$(aws --endpoint-url=$(ENDPOINT_URL) --profile $(PROFILE) sqs receive-message --queue-url $(QUEUE_URL) --max-number-of-messages 10 --output json); \
	  if [ "$$MESSAGES" == "" ]; then break; fi; \
	  for RECEIPT_HANDLE in $$(echo $$MESSAGES | jq -r '.Messages[].ReceiptHandle'); do \
	    echo "Deleting message: $$RECEIPT_HANDLE"; \
	    aws --endpoint-url=$(ENDPOINT_URL) --profile $(PROFILE) sqs delete-message --queue-url $(QUEUE_URL) --receipt-handle $$RECEIPT_HANDLE; \
	  done; \
	done
