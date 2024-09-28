SET TIMEZONE='UTC';
SET ORGANIZATION_ID  = '02fa7d30-c3de-4e0a-8f1e-2de120e7fd23';
SET CLUSTER_NAME = 'aws-jb-cirrus-load-testing-cluster';
SET START_DATE = '2024-09-26T00:00:00+00:00'::TIMESTAMP_TZ;
SET END_DATE   = '2024-09-30T00:00:00+00:00'::TIMESTAMP_TZ;
SET IMPOSSIBLE_VALUE = 100000000000000000000000;


WITH data AS (
    SELECT DISTINCT
     usage_date                       AS usage_date
    , organization_id                 AS organization_id
    , cloud_account_id                AS cloud_account_id
    , cluster_name                    AS cluster_name
    , LPAD(year(usage_date), 4, '0')  AS year
    , LPAD(month(usage_date), 2, '0') AS month
    , LPAD(day(usage_date), 2, '0')   AS day
    , LPAD(hour(usage_date), 2, '0')  AS hour
    , value                           AS value
    , labels                          AS labels
    FROM "PUBLIC"."PROMETHEUS_STAGING"
    WHERE 1=1
        AND usage_date::TIMESTAMP_TZ BETWEEN $START_DATE AND $END_DATE
        AND organization_id    = $ORGANIZATION_ID
        AND CLUSTER_NAME IN ($CLUSTER_NAME, $OTHER_CLUSTER_NAME)
)
, numeric_data AS (
    SELECT
      usage_date        AS usage_date
    , organization_id   AS organization_id
    , cloud_account_id  AS cloud_account_id
    , cluster_name      AS cluster_name
    , year              AS year
    , month             AS month
    , day               AS day
    , hour              AS hour
    , value             AS value
    , labels            AS labels
    FROM data
    -- Samples with values of NaN are not usable for numeric metrics.
    WHERE value::string != 'NaN'
    -- make sure we don't consume impossible values
    AND value < $IMPOSSIBLE_VALUE
)
, pod_working_set_memory AS (
    SELECT 
        DATE_TRUNC('hour', usage_date) AS metrics_hour,
        usage_date AS usage_date,
        cluster_name AS cluster_name,
        labels:node::string AS node_name,
        labels:namespace::string AS namespace,
        labels:pod::string AS kubernetes_pod_name,
        -- Calculate the percentage of node memory used by the pod for this interval
        value
    FROM numeric_data
    WHERE 1=1
        -- the following set is used to get POD level
        -- it includes all containers in the pod. It was discovered
        -- we may not get all containers in the pod when looking at the
        -- container set of working set memory metric, throwing off the calculation
        AND labels:__name__::string = 'container_memory_working_set_bytes'
        AND labels:node IS NOT NULL
        AND labels:pod IS NOT NULL
        AND labels:container IS NULL
        AND labels:image IS NULL
),
missing_records AS ( 
    SELECT 
        kubernetes_pod_name,
        usage_date,
        DATEDIFF(
            'second',
            LAG(usage_date) OVER (PARTITION BY kubernetes_pod_name ORDER BY usage_date),
            usage_date
        ) AS usage_time_diff_seconds,
        LAG(usage_date) OVER (PARTITION BY kubernetes_pod_name ORDER BY usage_date) AS previous_usage_date,
        cluster_name,
        node_name,
        namespace
    FROM pod_working_set_memory
)
SELECT 
    * 
FROM missing_records
WHERE usage_time_diff_seconds > 120
ORDER BY kubernetes_pod_name, usage_date DESC;