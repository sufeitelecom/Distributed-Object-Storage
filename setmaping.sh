#!/bin/bash

curl localhost:9200/metadata -XDELETE

curl localhost:9200/metadata -H"Content-Type: application/json" -XPUT -d'{"mappings":{"objects":{"properties":{"name":{"type":"keyword"},"version":{"type":"integer"},"size":{"type":"integer"},"hash":{"type":"keyword"}}}}}'