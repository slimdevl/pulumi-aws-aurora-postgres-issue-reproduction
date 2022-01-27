Issue Reproduction
==================

This repository reproduces a MAJOR issue with pulumi support for AWS Aurora (RDS) support. More specifically any calls to  `NewCluster` method thinks the `availabilityZones` change. Even if they are hardcoded.

So some state/diff logic is not doing the comparison correctly, because the AZs are static ... as one can observe by looking at the code.

I believe this is a **critical** issue making RDS/Aurora capabilities unusable. If someone uses this module, they will lose all their data each deployment.

There are two reproduction cases here:
1. Located in [quick-start-repro](quick-start-repro) is a reproduction of the issue for:  https://github.com/pulumi/pulumi-aws-quickstart-aurora-postgres. Follow the README there.
2. Located in [classic-repro](classic-repro) is a reproduction of the issue for:  https://github.com/pulumi/pulumi-aws. Follow the README there.

Because this happens with both QuickStart and Classic, I believe it may be rooted in the classic implementation, but I can't follow the code well enough to root cause it.