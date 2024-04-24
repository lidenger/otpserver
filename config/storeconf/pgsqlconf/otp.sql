/*
 Navicat Premium Data Transfer

 Source Server         : pgsql-local
 Source Server Type    : PostgreSQL
 Source Server Version : 160002 (160002)
 Source Host           : localhost:5432
 Source Catalog        : otp
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 160002 (160002)
 File Encoding         : 65001

 Date: 24/04/2024 12:13:45
*/


-- ----------------------------
-- Sequence structure for seq_otp_account_secret
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_otp_account_secret";
CREATE SEQUENCE "public"."seq_otp_account_secret" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_otp_operation_log
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_otp_operation_log";
CREATE SEQUENCE "public"."seq_otp_operation_log" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_otp_server
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_otp_server";
CREATE SEQUENCE "public"."seq_otp_server" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for seq_otp_server_ip_whitelist
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."seq_otp_server_ip_whitelist";
CREATE SEQUENCE "public"."seq_otp_server_ip_whitelist" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for otp_account_secret
-- ----------------------------
DROP TABLE IF EXISTS "public"."otp_account_secret";
CREATE TABLE "public"."otp_account_secret" (
  "id" int8 NOT NULL DEFAULT nextval('seq_otp_account_secret'::regclass),
  "secret_seed_cipher" varchar(500) COLLATE "pg_catalog"."default",
  "account" varchar(100) COLLATE "pg_catalog"."default",
  "is_enable" int2,
  "data_check" varchar(100) COLLATE "pg_catalog"."default",
  "create_time" timestamp(6),
  "update_time" timestamp(6)
)
;
COMMENT ON COLUMN "public"."otp_account_secret"."id" IS '主键自增ID';
COMMENT ON COLUMN "public"."otp_account_secret"."secret_seed_cipher" IS '密钥种子密文 = AES(KEY, 密钥种子)';
COMMENT ON COLUMN "public"."otp_account_secret"."account" IS '账号';
COMMENT ON COLUMN "public"."otp_account_secret"."is_enable" IS '是否启用，1启用，2禁用';
COMMENT ON COLUMN "public"."otp_account_secret"."data_check" IS '数据校验 = HMACSHA256(KEY, secret_seed_cipher + account + is_enable)';
COMMENT ON COLUMN "public"."otp_account_secret"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."otp_account_secret"."update_time" IS '更新时间';
COMMENT ON TABLE "public"."otp_account_secret" IS '账号密钥';

-- ----------------------------
-- Table structure for otp_operation_log
-- ----------------------------
DROP TABLE IF EXISTS "public"."otp_operation_log";
CREATE TABLE "public"."otp_operation_log" (
  "id" int8 NOT NULL DEFAULT nextval('seq_otp_operation_log'::regclass),
  "item" varchar(200) COLLATE "pg_catalog"."default",
  "content" varchar(800) COLLATE "pg_catalog"."default",
  "ip" varchar(100) COLLATE "pg_catalog"."default",
  "create_time" timestamp(6),
  "update_time" timestamp(6)
)
;
COMMENT ON COLUMN "public"."otp_operation_log"."id" IS '主键自增ID';
COMMENT ON COLUMN "public"."otp_operation_log"."item" IS '操作项';
COMMENT ON COLUMN "public"."otp_operation_log"."content" IS '操作内容';
COMMENT ON COLUMN "public"."otp_operation_log"."ip" IS '操作者IP';
COMMENT ON COLUMN "public"."otp_operation_log"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."otp_operation_log"."update_time" IS '更新时间';
COMMENT ON TABLE "public"."otp_operation_log" IS '操作日志';

-- ----------------------------
-- Table structure for otp_server
-- ----------------------------
DROP TABLE IF EXISTS "public"."otp_server";
CREATE TABLE "public"."otp_server" (
  "id" int8 NOT NULL DEFAULT nextval('seq_otp_account_secret'::regclass),
  "server_sign" varchar(100) COLLATE "pg_catalog"."default",
  "server_name" varchar(200) COLLATE "pg_catalog"."default",
  "server_secret_cipher" varchar(500) COLLATE "pg_catalog"."default",
  "server_secret_iv" varchar(500) COLLATE "pg_catalog"."default",
  "server_remark" varchar(300) COLLATE "pg_catalog"."default",
  "is_enable" int2,
  "is_operate_sensitive_data" int2,
  "is_enable_iplist" int2,
  "data_check" varchar(100) COLLATE "pg_catalog"."default",
  "create_time" timestamp(6),
  "update_time" timestamp(6)
)
;
COMMENT ON COLUMN "public"."otp_server"."id" IS '主键自增ID';
COMMENT ON COLUMN "public"."otp_server"."server_sign" IS '服务标识';
COMMENT ON COLUMN "public"."otp_server"."server_name" IS '服务名称';
COMMENT ON COLUMN "public"."otp_server"."server_secret_cipher" IS '服务密钥';
COMMENT ON COLUMN "public"."otp_server"."server_secret_iv" IS '服务密钥IV';
COMMENT ON COLUMN "public"."otp_server"."server_remark" IS '服务描述';
COMMENT ON COLUMN "public"."otp_server"."is_enable" IS '是否启用，1启用，2禁用';
COMMENT ON COLUMN "public"."otp_server"."is_operate_sensitive_data" IS '是否可以操作敏感数据（例如：密钥数据），1是，2否';
COMMENT ON COLUMN "public"."otp_server"."is_enable_iplist" IS '是否启用服务IP白名单，1启用，2禁用';
COMMENT ON COLUMN "public"."otp_server"."data_check" IS '数据校验 ';
COMMENT ON COLUMN "public"."otp_server"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."otp_server"."update_time" IS '更新时间';
COMMENT ON TABLE "public"."otp_server" IS '接入服务';

-- ----------------------------
-- Table structure for otp_server_ip_whitelist
-- ----------------------------
DROP TABLE IF EXISTS "public"."otp_server_ip_whitelist";
CREATE TABLE "public"."otp_server_ip_whitelist" (
  "id" int8 NOT NULL DEFAULT nextval('seq_otp_server_ip_whitelist'::regclass),
  "ip" varchar(100) COLLATE "pg_catalog"."default",
  "server_sign" varchar(100) COLLATE "pg_catalog"."default",
  "create_time" timestamp(6),
  "update_time" timestamp(6)
)
;
COMMENT ON COLUMN "public"."otp_server_ip_whitelist"."id" IS '主键自增ID';
COMMENT ON COLUMN "public"."otp_server_ip_whitelist"."ip" IS 'ip';
COMMENT ON COLUMN "public"."otp_server_ip_whitelist"."server_sign" IS '服务标识';
COMMENT ON COLUMN "public"."otp_server_ip_whitelist"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."otp_server_ip_whitelist"."update_time" IS '更新时间';
COMMENT ON TABLE "public"."otp_server_ip_whitelist" IS '接入服务IP白名单';

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."seq_otp_account_secret"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."seq_otp_operation_log"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."seq_otp_server"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."seq_otp_server_ip_whitelist"', 1, false);

-- ----------------------------
-- Indexes structure for table otp_account_secret
-- ----------------------------
CREATE INDEX "idx_otp_account_secret_account" ON "public"."otp_account_secret" USING btree (
  "account" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_otp_account_secret_update_time" ON "public"."otp_account_secret" USING btree (
  "update_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table otp_account_secret
-- ----------------------------
ALTER TABLE "public"."otp_account_secret" ADD CONSTRAINT "otp_account_secret_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table otp_operation_log
-- ----------------------------
CREATE INDEX "idx_otp_operation_log_update_time" ON "public"."otp_operation_log" USING btree (
  "update_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table otp_operation_log
-- ----------------------------
ALTER TABLE "public"."otp_operation_log" ADD CONSTRAINT "otp_operation_log_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table otp_server_ip_whitelist
-- ----------------------------
CREATE INDEX "idx_otp_server_ip_whitelist_ip" ON "public"."otp_server_ip_whitelist" USING btree (
  "ip" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_otp_server_ip_whitelist_server_sign" ON "public"."otp_server_ip_whitelist" USING btree (
  "update_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE INDEX "idx_otp_server_ip_whitelist_update_time" ON "public"."otp_server_ip_whitelist" USING btree (
  "update_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table otp_server_ip_whitelist
-- ----------------------------
ALTER TABLE "public"."otp_server_ip_whitelist" ADD CONSTRAINT "otp_server_ip_whitelist_pkey" PRIMARY KEY ("id");
