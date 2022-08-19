resource "alkira_byoip_prefix" "test" {
  prefix      = "172.16.1.2"
  cxp         = "US-WEST-1"
  description = "simple test"
  message     = "1|aws|0123456789AB|198.51.100.0/24|20211231|SHA256|RSAPSS"
  signature   = "rDr1314-SVB--C~77Lb6xRH1m9kcELJvmtmci1~Qzt4dW9cE0a6y33KOQXFPEW4AKdgyPxEnfhLGoZLurGZWBAfZrhiVHc7CMo3n2GEFKewgdBTWKjS-1UWyHzLjg9yD-zCUZl08AlrUjLIh-jBz~n7RbQGULz8UpSrKkCvsAN-U1gAg5xuyeBhcQkN7djjQo8eYSLvxAzp~TlzuKiDt~3tbiyQ0gC8US-bvAZno6-PDS64RABc3Z-6~C~B4SJQmvAPlJq460xly-VjJWPnSyxjbi~qT4O4eTfC4~SaTMGJuCJ3FZivO2YhGW-UBUn5Io7R8x8~3mhzOn3VWtUDQUQ__"
  public_key  = "-----BEGIN PUBLIC KEY-----MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzQzZ+Z6Xlu6PQ8n9uN06nGih2PSNbOpeAo8e5U7DObKIjgow6JODnZSFrmp9tIWfqYlQ87msuxBXQDrDq0lFU3UYay4quEOh6vn+zChPGUVfkgp7Nm1m3wvD8dx0Zdn/AYoXZa06umg2RiGQwSQlld9JNy+L2K7QzM/lVIOj659ft7boBjnhNzhrbqjshISsCn9z5KKuMnUqof5Vd7RPGZskYnYDscLFV0dKOuFCFa2wW57NPqLKuCukcfIQUWQkhLLTi3tlRQd0X4ElfpF7gxJEQPY+9FanwaA3Yp3hEj5YX9z0Nz0YkTrnCnhMxINxxPcLAnca1SmwvVt+9uuf2QIDAQAB-----END PUBLIC KEY-----"
}
