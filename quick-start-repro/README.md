Quick Start Issue Reproduction
=============================

As the log below shows, on the second run of this... The `NewCluster` method thinks the `availabilityZones` change. So some state/diff logic is not doing the comparison correctly, because the AZs are static ... as one can observe by looking at the code.

I believe this is a critical issue making RDS/Aurora capabilities unusable. If someone uses this module, they will lose all their data each deployment.


# Reproduction using the quick start
1. run `export AWS_REGION=us-east-1`
2. edit main.go to use your email.
3. run `pulumi stack init repro`
4. run `pulumi up -y`
5. Do it again to see the issue `pulumi up -y`

# Expected output
No change in resources on any deployment after the first run.

# Actual output
The database is always recreated, which causes data loss.

## Logs
```bash
pulumi up -y
Previewing update (repro)

View Live: https://app.pulumi.com/slimdevl/simple-postgres-go/repro/previews/9cc4ba77-6294-4460-9300-ac2c0b16fef9

     Type                                             Name                                                         Plan
 +   pulumi:pulumi:Stack                              simple-postgres-go-repro                                     create
 +   ├─ aws-quickstart-vpc:index:Vpc                  simple-vpc                                                   create
 +   │  ├─ aws:cloudwatch:LogGroup                    simple-vpc-flow-logs                                         create
 +   │  ├─ aws:iam:Role                               simple-vpc-vpc-flow-log-role                                 create
 +   │  ├─ aws:ec2:Vpc                                simple-vpc-vpc                                               create
 +   │  ├─ aws:iam:RolePolicy                         simple-vpc-vpc-flow-log-policy                               create
 +   │  ├─ aws:ec2:Subnet                             simple-vpc-private-subnet-a-0                                create
 +   │  ├─ aws:ec2:Subnet                             simple-vpc-private-subnet-a-1                                create
 +   │  ├─ aws:ec2:FlowLog                            simple-vpc-vpc-flow-log                                      create
 +   │  └─ aws:ec2:InternetGateway                    simple-vpc-internet-gateway                                  create
 +   └─ aws-quickstart-aurora-postgres:index:Cluster  smiple-aurora-postgres                                       create
 +      ├─ aws:kms:Key                                smiple-aurora-postgres-database-kms-key                      create
 +      ├─ aws:rds:SubnetGroup                        smiple-aurora-postgres-db-subnet-group                       create
 +      ├─ aws:rds:ParameterGroup                     smiple-aurora-postgres-db-parameter-group                    create
 +      ├─ aws:rds:ClusterParameterGroup              smiple-aurora-postgres-parameter-group                       create
 +      ├─ aws:sns:Topic                              sns-topic                                                    create
 +      ├─ aws:kms:Alias                              smiple-aurora-postgres-database-kms-key-alias                create
 +      ├─ aws:rds:Cluster                            smiple-aurora-postgres-postgresql-cluster                    create
 +      │  └─ aws:rds:ClusterInstance                 smiple-aurora-postgres-aurora-database-0                     create
 +      │     ├─ aws:rds:EventSubscription            smiple-aurora-postgres-parameter-group-event-subscription-0  create
 +      │     ├─ aws:rds:EventSubscription            smiple-aurora-postgres-cluster-event-subscription-0          create
 +      │     ├─ aws:rds:EventSubscription            smiple-aurora-postgres-instance-event-subscription-0         create
 +      │     ├─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-cpu-alarm-0                           create
 +      │     ├─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-max-used-tx-alarm-0                   create
 +      │     └─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-free-local-storage-alarm-0            create
 +      └─ aws:sns:TopicSubscription                  sns-topic-subscription                                       create

Resources:
    + 26 to create

Updating (repro)

View Live: https://app.pulumi.com/slimdevl/simple-postgres-go/repro/updates/3

     Type                                             Name                                                         Status
 +   pulumi:pulumi:Stack                              simple-postgres-go-repro                                     created
 +   ├─ aws-quickstart-vpc:index:Vpc                  simple-vpc                                                   created
 +   │  ├─ aws:ec2:Vpc                                simple-vpc-vpc                                               created
 +   │  ├─ aws:cloudwatch:LogGroup                    simple-vpc-flow-logs                                         created
 +   │  ├─ aws:iam:Role                               simple-vpc-vpc-flow-log-role                                 created
 +   │  ├─ aws:iam:RolePolicy                         simple-vpc-vpc-flow-log-policy                               created
 +   │  ├─ aws:ec2:InternetGateway                    simple-vpc-internet-gateway                                  created
 +   │  ├─ aws:ec2:Subnet                             simple-vpc-private-subnet-a-1                                created
 +   │  ├─ aws:ec2:FlowLog                            simple-vpc-vpc-flow-log                                      created
 +   │  └─ aws:ec2:Subnet                             simple-vpc-private-subnet-a-0                                created
 +   └─ aws-quickstart-aurora-postgres:index:Cluster  smiple-aurora-postgres                                       created
 +      ├─ aws:sns:Topic                              sns-topic                                                    created
 +      ├─ aws:rds:ParameterGroup                     smiple-aurora-postgres-db-parameter-group                    created
 +      ├─ aws:rds:ClusterParameterGroup              smiple-aurora-postgres-parameter-group                       created
 +      ├─ aws:rds:SubnetGroup                        smiple-aurora-postgres-db-subnet-group                       created
 +      ├─ aws:kms:Key                                smiple-aurora-postgres-database-kms-key                      created
 +      ├─ aws:sns:TopicSubscription                  sns-topic-subscription                                       created
 +      ├─ aws:kms:Alias                              smiple-aurora-postgres-database-kms-key-alias                created
 +      └─ aws:rds:Cluster                            smiple-aurora-postgres-postgresql-cluster                    created
 +         └─ aws:rds:ClusterInstance                 smiple-aurora-postgres-aurora-database-0                     created
 +            ├─ aws:rds:EventSubscription            smiple-aurora-postgres-cluster-event-subscription-0          created
 +            ├─ aws:rds:EventSubscription            smiple-aurora-postgres-instance-event-subscription-0         created
 +            ├─ aws:rds:EventSubscription            smiple-aurora-postgres-parameter-group-event-subscription-0  created
 +            ├─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-max-used-tx-alarm-0                   created
 +            ├─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-cpu-alarm-0                           created
 +            └─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-free-local-storage-alarm-0            created

Resources:
    + 26 created

Duration: 11m10s

❯ pulumi up -y
Previewing update (repro)

View Live: https://app.pulumi.com/slimdevl/simple-postgres-go/repro/previews/52a89561-2372-4382-836d-e0495301f713

     Type                                             Name                                                  Plan        Info
     pulumi:pulumi:Stack                              simple-postgres-go-repro
     └─ aws-quickstart-aurora-postgres:index:Cluster  smiple-aurora-postgres
 +-     └─ aws:rds:Cluster                            smiple-aurora-postgres-postgresql-cluster             replace     [diff: ~availabilityZones]
 +-        └─ aws:rds:ClusterInstance                 smiple-aurora-postgres-aurora-database-0              replace     [diff: ~clusterIdentifier]
 ~            ├─ aws:rds:EventSubscription            smiple-aurora-postgres-cluster-event-subscription-0   update      [diff: ~sourceIds]
 ~            ├─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-free-local-storage-alarm-0     update      [diff: ~dimensions]
 ~            ├─ aws:rds:EventSubscription            smiple-aurora-postgres-instance-event-subscription-0  update      [diff: ~sourceIds]
 ~            ├─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-max-used-tx-alarm-0            update      [diff: ~dimensions]
 ~            └─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-cpu-alarm-0                    update      [diff: ~dimensions]

Resources:
    ~ 5 to update
    +-2 to replace
    7 changes. 19 unchanged

Updating (repro)

View Live: https://app.pulumi.com/slimdevl/simple-postgres-go/repro/updates/4

     Type                                             Name                                       Status                   Info
     Type                                             Name                                                  Status       Info
     pulumi:pulumi:Stack                              simple-postgres-go-repro
     └─ aws-quickstart-aurora-postgres:index:Cluster  smiple-aurora-postgres
 +-     └─ aws:rds:Cluster                            smiple-aurora-postgres-postgresql-cluster             replaced     [diff: ~availabilityZones]
 +-        └─ aws:rds:ClusterInstance                 smiple-aurora-postgres-aurora-database-0              replaced     [diff: ~clusterIdentifier]
 ~            ├─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-cpu-alarm-0                    updated      [diff: ~dimensions]
 ~            ├─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-max-used-tx-alarm-0            updated      [diff: ~dimensions]
 ~            ├─ aws:rds:EventSubscription            smiple-aurora-postgres-instance-event-subscription-0  updated      [diff: ~sourceIds]
 ~            ├─ aws:cloudwatch:MetricAlarm           smiple-aurora-postgres-free-local-storage-alarm-0     updated      [diff: ~dimensions]
 ~            └─ aws:rds:EventSubscription            smiple-aurora-postgres-cluster-event-subscription-0   updated      [diff: ~sourceIds]

Resources:
    ~ 5 updated
    +-2 replaced
    7 changes. 19 unchanged

```