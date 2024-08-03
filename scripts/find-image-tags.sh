for i in $(crane ls ghcr.io/cloudzero/cloudzero-insights-controller/cloudzero-insights-controller | grep -E '^[0-9]+\.[0-9]+\.[0-9]+$|^develop$|^main$|^latest$'); do
  echo "TAG: ${i}"
done
