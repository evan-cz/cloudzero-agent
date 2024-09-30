use database live_billing;
SET ORGANIZATION_ID  = '02fa7d30-c3de-4e0a-8f1e-2de120e7fd23';
SET CLUSTER_NAME = 'cloudzero-eks-cluster-eksCluster-45e897d';
SET YEAR_TIME   = '2024';
SET MONTH_TIME  = '09';
SET DAY_TIME    = '30';

WITH TIMESERIES AS (
    SELECT
        SEQ8() AS RECORD_ID,
        cloud_account_id,
        day,
        month,
        year,
        organization_id,
        cluster_name,
        hour,
        FLATTENED.VALUE AS DATA
    FROM
        "PUBLIC"."prometheus_container_external_org_02fa7d30_c3de_4e0a_8f1e_2de120e7fd23",
        LATERAL FLATTEN(INPUT => VALUE:timeseries) FLATTENED
    WHERE
        organization_id = $ORGANIZATION_ID
        AND cluster_name = $CLUSTER_NAME
        AND year        = $YEAR_TIME
        AND month       = $MONTH_TIME
        AND day         in ($DAY_TIME)
        AND hour        in ('00', '01', '02')
)
-- select * from TIMESERIES;
, DATA_ROWS AS (
    SELECT DISTINCT
        MAX(TO_TIMESTAMP_TZ(DATA:samples[0]:timestamp::STRING)) AS USAGE_DATE,
        MAX(organization_id)::STRING        AS organization_id,
        MAX(cloud_account_id)::STRING       AS cloud_account_id,
        MAX(cluster_name)::STRING           AS cluster_name,
        MAX(year)::STRING                   AS year,
        MAX(month)::STRING                  AS month,
        MAX(day)::STRING                    AS day,
        MAX(hour)::STRING                   AS hour,
        MAX(DATA:samples[0]:value::FLOAT)   AS VALUE,
        OBJECT_AGG(VALUE:name, VALUE:value) AS LABELS
    FROM
        TIMESERIES,
        LATERAL FLATTEN(INPUT => DATA:labels)
    GROUP BY
        RECORD_ID
)
, pod_working_set_memory AS (
    SELECT
        DATE_TRUNC('hour', usage_date)   AS metrics_hour,
        USAGE_DATE                       AS usage_date,
        organization_id                  AS organization_id,
        cluster_name                     AS cluster_name,
        labels:node::string              AS node_name,
        labels:namespace::string         AS namespace,
        labels:pod::string               AS kubernetes_pod_name,
        value::FLOAT                     AS value,
        labels                           AS labels
    FROM DATA_ROWS
    WHERE 1=1
        AND labels:__name__::string = 'container_memory_working_set_bytes'
        AND labels:node IS NOT NULL
        AND labels:pod IS NOT NULL
        AND labels:container IS NULL
        AND labels:image IS NULL
)
-- SELECT
--     metrics_hour,
--     usage_date,
--     organization_id,
--     cluster_name,
--     node_name,
--     namespace,
--     kubernetes_pod_name,
--     value,
--     labels
-- FROM pod_working_set_memory
-- ORDER BY
--     cluster_name,node_name,namespace,kubernetes_pod_name,usage_date desc;
, missing_records AS ( 
    SELECT
        DATE_TRUNC('hour', usage_date) AS metrics_hour,
        kubernetes_pod_name            AS kubernetes_pod_name,
        usage_date                     AS usage_date,
        DATEDIFF(
            'second',
            LAG(usage_date) OVER (PARTITION BY cluster_name, node_name, namespace, kubernetes_pod_name ORDER BY usage_date),
            usage_date
        )                              AS usage_time_diff_seconds,
        LAG(usage_date) OVER (PARTITION BY cluster_name, node_name, namespace, kubernetes_pod_name ORDER BY usage_date) AS previous_usage_date,
        cluster_name                   AS cluster_name,
        node_name                      AS node_name,
        namespace                      AS namespace,
        value                          AS value
    FROM pod_working_set_memory
)
SELECT 
    metrics_hour,
    usage_time_diff_seconds
    usage_date,
    previous_usage_date,
    namespace,
    kubernetes_pod_name,
    value,
    cluster_name,
    node_name
FROM missing_records
--WHERE usage_time_diff_seconds > 120
ORDER BY cluster_name, namespace, kubernetes_pod_name, usage_date DESC;

