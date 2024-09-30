SET TIMEZONE='UTC';
SET ORGANIZATION_ID  = '02fa7d30-c3de-4e0a-8f1e-2de120e7fd23';
SET CLUSTER_NAME = 'cloudzero-eks-cluster-eksCluster-45e897d';
SET START_DATE = '2024-09-29T00:00:00+00:00'::TIMESTAMP_TZ;
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
        AND CLUSTER_NAME IN ($CLUSTER_NAME)
)
, node_info AS (
    -- find the unique pod info records
    SELECT
        DISTINCT 
        usage_date AS usage_date,
        cluster_name AS cluster_name,
        labels:node::string AS node_name,
        value AS value,
        labels AS labels
    FROM data
    WHERE labels:__name__::string = 'kube_node_info'
    -- You can omit GROUP BY if using DISTINCT
)
, missing_node_records AS ( 
    SELECT 
        usage_date,
        cluster_name,
        node_name,
        DATEDIFF(
            'second',
            LAG(usage_date) OVER (PARTITION BY cluster_name, node_name ORDER BY usage_date),
            usage_date
        ) AS usage_time_diff_seconds,
        LAG(usage_date) OVER (PARTITION BY cluster_name, node_name ORDER BY usage_date) AS previous_usage_date
    FROM node_info 
)
SELECT 
    * 
FROM missing_node_records
-- over 15 minutes between records for the same node
--WHERE usage_time_diff_seconds > 60
ORDER BY cluster_name, node_name, usage_date asc;