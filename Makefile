# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

DIR = $(shell pwd)/app

CONFIG_PATH = $(shell pwd)/config
IDL_PATH = $(shell pwd)/idl

SERVICES := gateway user favorite search_engine
service = $(word 1, $@)

node = 0

BIN = $(shell pwd)/bin

.PHONY: proto
proto:
	@for file in $(IDL_PATH)/*.proto; do \
		protoc -I $(IDL_PATH) $$file --go-grpc_out=$(IDL_PATH)/pb --go_out=$(IDL_PATH)/pb; \
	done
	@for file in $(shell find $(IDL_PATH)/pb/* -type f); do \
		protoc-go-inject-tag -input=$$file; \
	done

.PHONY: pyproto
pyproto:
	python3 -m grpc_tools.protoc --grpc_python_out=$(IDL_PATH)/pb/search_vector \
		--python_out=$(IDL_PATH)/pb/search_vector --pyi_out=$(IDL_PATH)/pb/search_vector \
		--proto_path=$(IDL_PATH) $(IDL_PATH)/search_vector.proto;

# python3 -m grpc_tools.protoc -I ./ --python_out=./ --grpc_python_out=. ./search_vector.proto


.PHONY: $(SERVICES)
$(SERVICES):
	go build -o $(BIN)/$(service) $(DIR)/$(service)/cmd
	$(BIN)/$(service) -config $(CONFIG_PATH) -srvnum=$(node)

.PHONY: env-up
env-up:
	docker-compose -f docker-compose-milvus.yaml up -d
	docker-compose -f docker-compose.yaml up -d
	docker-compose -f docker-compose-with-kafka.yaml up -d

.PHONY: env-down
env-down:
	docker-compose down

.PHONY: run-%
run-%:
	go run $(DIR)/$*/cmd/main.go;

.PHONY: python-start
python-start:
    # export PROTOCOL_BUFFERS_PYTHON_IMPLEMENTATION=python
	python3 main.py

.PHONY: python-consume
python-consume:
	python3 ./vector_index.py
