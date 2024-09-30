use database live_billing;
SET ORGANIZATION_ID  = '02fa7d30-c3de-4e0a-8f1e-2de120e7fd23';
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
        AND year        = $YEAR_TIME
        AND month       = $MONTH_TIME
        AND day         in ($DAY_TIME)
        AND hour        in ('00', '01', '02', '03', '04', '05', '06', '07', '08', '09', '10', '11', '12', '13', '14', '15', '16', '17', '18', '19', '20', '21', '22', '23')
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
, FLAT_DATA AS (
    SELECT
        DATE_TRUNC('hour', usage_date)   AS metrics_hour,
        USAGE_DATE                       AS usage_date,
        organization_id                  AS organization_id,
        cluster_name                     AS cluster_name,
        value::FLOAT                     AS value,
        labels                           AS labels
    FROM DATA_ROWS
    WHERE labels:"__name__"::STRING = 'kube_pod_labels'
)
, LABELS_FLATTENED AS (
    SELECT
        BD.metrics_hour      AS metrics_hour
        , lf.key             AS label_key
    FROM FLAT_DATA BD,
        LATERAL FLATTEN(input => BD.LABELS) lf
    WHERE 1=1
        AND lf.key LIKE 'label_%'
)
, UNIQUE_LABEL_KEYS_PER_DAY AS (
    SELECT
        metrics_hour                                           AS metrics_hour
        , ARRAY_AGG(DISTINCT REPLACE(label_key, 'label_', '')) AS unique_label_keys_unsorted
        , COUNT(DISTINCT label_key)                            AS unique_label_keys_count
    FROM LABELS_FLATTENED
    GROUP BY
        metrics_hour
)
, SORTED_LABEL_KEYS AS (
    SELECT
        metrics_hour                             AS metrics_hour
        , ARRAY_SORT(unique_label_keys_unsorted) AS unique_label_keys_list
        , unique_label_keys_count
    FROM UNIQUE_LABEL_KEYS_PER_DAY
)
SELECT
    metrics_hour,
    unique_label_keys_count,
    unique_label_keys_list
FROM
    sorted_label_keys
ORDER BY
    metrics_hour;
