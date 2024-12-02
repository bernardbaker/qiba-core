#!/bin/bash
gcloud config set project qiba-core-441819

for i in {0..49}
do
  gcloud compute addresses create qiba-ip-"$i" \
    --region=us-central1 \
    --subnet=qiba
done

for i in {0..49}
do
  if [ "$(gcloud compute addresses describe qiba-ip-"$i" --region=us-central1 --format="value(status)")" != "RESERVED" ]; then
    echo "qiba-ip-$i is not RESERVED";
    exit 1;
  fi
done

for i in {0..49}
do
  gcloud compute forwarding-rules create qiba-"$i" \
    --region=us-central1 \
    --network=default \
    --address=qiba-ip-"$i" \
    --allow-psc-global-access \
    --target-service-attachment=projects/p-xsl3fefiaarezolovqxmlbwv/regions/us-central1/serviceAttachments/sa-us-central1-67472b9fff0ef50e09e322d4-"$i"
done

if [ "$(gcloud compute forwarding-rules list --regions=us-central1 --format="csv[no-heading](name)" --filter="(name:qiba*)" | wc -l)" -gt 50 ]; then
  echo "Project has too many forwarding rules that match prefix qiba. Either delete the competing resources or choose another endpoint prefix."
  exit 2;
fi

gcloud compute forwarding-rules list --regions=us-central1 --format="json(IPAddress,name)" --filter="name:(qiba*)" > atlasEndpoints-qiba.json
