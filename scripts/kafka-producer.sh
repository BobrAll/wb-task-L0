#!/bin/bash

TOPIC="orders"
BROKER="localhost:9092"

BASE_CHRT_ID=9934930

for i in {1..15}
do
  ORDER_UID="b563feb7b2b84b6test$i"
  TRACK_NUMBER="WBILMTESTTRACK$i"
  TRANSACTION="b563feb7b2b84b6test$i"
  CUSTOMER_ID="test$i"
  DATE_CREATED="2025-08-24T02:16:$(printf "%02d" $i)Z"
  CHRT_ID=$((BASE_CHRT_ID + i))

  MESSAGE=$(cat <<EOF
{"order_uid":"$ORDER_UID","track_number":"$TRACK_NUMBER","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"$TRANSACTION","request_id":"","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":$CHRT_ID,"track_number":"$TRACK_NUMBER","price":453,"rid":"ab4219087a764ae0btest$i","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"$CUSTOMER_ID","delivery_service":"meest","shardkey":"9","sm_id":99,"date_created":"$DATE_CREATED","oof_shard":"1"}
EOF
)

  echo "Sending order #$i with chrt_id=$CHRT_ID"
  echo "$MESSAGE" | kafka-console-producer --bootstrap-server "$BROKER" --topic "$TOPIC"
done
