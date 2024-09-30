SET TIMEZONE='UTC';
SET ORGANIZATION_ID  = '02fa7d30-c3de-4e0a-8f1e-2de120e7fd23';
SET CLUSTER_NAME = 'cloudzero-eks-cluster-eksCluster-45e897d';
SET START_DATE = '2024-09-30T00:00:00+00:00'::TIMESTAMP_TZ;
SET END_DATE   = '2024-09-30T02:00:00+00:00'::TIMESTAMP_TZ;
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
        AND CLUSTER_NAME IN ($CLUSTER_NAME)
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
, container_cpu_usage_seconds AS (
    SELECT
        DATE_TRUNC('hour', usage_date) AS metrics_hour,
        cluster_name                   AS cluster_name,
        labels:node::string            AS node_name,
        labels:namespace::string       AS namespace,
        labels:pod::string             AS kubernetes_pod_name,
        labels:container::string       AS container,
        usage_date                     AS usage_date,
        value                          AS cpu_usage_seconds_total,
        LAG(value) OVER (PARTITION BY cluster_name, node_name, kubernetes_pod_name, container ORDER BY usage_date) AS previous_cpu_usage_total,
        LAG(usage_date) OVER (PARTITION BY cluster_name, node_name, kubernetes_pod_name, container ORDER BY usage_date) AS previous_usage_date
    FROM
        numeric_data
    WHERE 
        labels:__name__::string = 'container_cpu_usage_seconds_total'
        AND labels:image IS NOT NULL
        AND labels:container IS NOT NULL
)
-- 
-- select 
--     metrics_hour,
--     count(distinct node_name) as unique_nodes
-- from container_cpu_usage_seconds
-- group by metrics_hour;

, missing_records AS ( 
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
    FROM container_cpu_usage_seconds
)
SELECT 
    * 
FROM missing_records
-- WHERE usage_time_diff_seconds > 120
ORDER BY kubernetes_pod_name, usage_date asc;