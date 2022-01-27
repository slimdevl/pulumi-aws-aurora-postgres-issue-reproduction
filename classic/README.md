Quick Start Issue Reproduction
=============================

As the log below shows, on the second run of this... The `NewCluster` method thinks the `availabilityZones` change. So some state/diff logic is not doing the comparison correctly, because the AZs are static ... as one can observe by looking at the code.

I believe this is a critical issue making RDS/Aurora capabilities unusable. If someone uses this module, they will lose all their data each deployment.

# Reproduction using the quick start
1. run `export AWS_REGION=us-east-1`
2. run `pulumi stack init repro`
3. run `pulumi up -y`
4. Do it again to see the issue `pulumi up -y`

# Expected output
No change in resources on any deployment after the first run.

# Actual output
The database is always recreated, which causes data loss.

## Logs (Issue 1)
```bash
❯ pulumi up -y
Previewing update (repro)

View Live: https://app.pulumi.com/slimdevl/simple-aws-classic-postgres-go/repro/previews/21266701-525d-49d8-ae4b-e194fa92baec

     Type                        Name                                  Plan        Info
     pulumi:pulumi:Stack         simple-aws-classic-postgres-go-repro
 +-  ├─ aws:rds:Cluster          aws-classic-postgres-db-cluster       replace     [diff: ~availabilityZones]
 +-  ├─ aws:rds:ClusterInstance  aws-classic-cluster-ins-us-east-1b    replace     [diff: ~clusterIdentifier]
 +-  └─ aws:rds:ClusterInstance  aws-classic-cluster-ins-us-east-1a    replace     [diff: ~clusterIdentifier]

Resources:
    +-3 to replace
    3 unchanged

Updating (repro)

View Live: https://app.pulumi.com/slimdevl/simple-aws-classic-postgres-go/repro/updates/2

     Type                 Name                                  Status                   Info
     pulumi:pulumi:Stack  simple-aws-classic-postgres-go-repro  **failed**               1 error
 +-  └─ aws:rds:Cluster   aws-classic-postgres-db-cluster       **replacing failed**     [diff: ~availabilityZones]; 1 error

Diagnostics:
  pulumi:pulumi:Stack (simple-aws-classic-postgres-go-repro):
    error: update failed

  aws:rds:Cluster (aws-classic-postgres-db-cluster):
    error: 1 error occurred:
    	* error creating RDS cluster: DBClusterAlreadyExistsFault: DB Cluster already exists
    	status code: 400, request id: 32da8d6f-6c43-477a-bc4a-da0cee8fd153

Outputs:
    subnets: [
        [0]: "subnet-2a46fc4c"
        [1]: "subnet-72f77753"
        [2]: "subnet-143abd4b"
        [3]: "subnet-66c0bc68"
        [4]: "subnet-72f77753"
        [5]: "subnet-d4c40de5"
        [6]: "subnet-2a46fc4c"
        [7]: "subnet-aa0653e7"
    ]

Resources:
    3 unchanged

Duration: 4s
```