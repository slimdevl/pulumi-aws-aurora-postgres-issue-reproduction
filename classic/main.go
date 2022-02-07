package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		t := true
		vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{Default: &t})
		if err != nil {
			return err
		}
		subnetIds, err := ec2.GetSubnetIds(ctx, &ec2.GetSubnetIdsArgs{VpcId: vpc.Id})
		if err != nil {
			return err
		}
		// Hardcode so we know the order doesn't changes
		availabilityZones := []string{"us-east-1a", "us-east-1b"}
		subnets := make([]string, 2)
		for _, id := range subnetIds.Ids {
			selected, err := ec2.LookupSubnet(ctx, &ec2.LookupSubnetArgs{
				Id: &id,
			}, nil)
			if err != nil {
				return err
			}
			if selected.DefaultForAz { // force alignment with indexes of AZ
				if selected.AvailabilityZone == "us-east-1a" {
					subnets[0] = selected.Id
				} else if selected.AvailabilityZone == "us-east-1b" {
					subnets[1] = selected.Id
				}
			}
			// ctx.Export(fmt.Sprintf("subnet-%d.DefaultForAz", i), pulumi.Bool(selected.DefaultForAz))
			// ctx.Export(fmt.Sprintf("subnet-%d.AvailabilityZone", i), pulumi.String(selected.AvailabilityZone))
		}
		ctx.Export("subnets", pulumi.ToStringArray(subnets))

		subnetGroup, err := rds.NewSubnetGroup(ctx,
			"aws-classic-aurora-postgres-subnet-group",
			&rds.SubnetGroupArgs{
				SubnetIds: pulumi.StringArray{
					pulumi.String(subnets[0]),
					pulumi.String(subnets[1]),
				},
				Tags: pulumi.StringMap{
					"Name": pulumi.String("aurora postgres subnet group"),
				},
			})
		if err != nil {
			return err
		}
		pgInst, err := rds.NewParameterGroup(
			ctx, "aws-classic-aurora-postgresql12",
			&rds.ParameterGroupArgs{
				Description: pulumi.StringPtr("aurora postgres 12 parameter groupu"),
				Family:      pulumi.String("aurora-postgresql12"),
				Parameters: rds.ParameterGroupParameterArray{
					&rds.ParameterGroupParameterArgs{
						Name:  pulumi.String("log_rotation_age"),
						Value: pulumi.String("1440"),
					},
					&rds.ParameterGroupParameterArgs{
						Name:  pulumi.String("log_rotation_size"),
						Value: pulumi.String("102400"),
					},
				},
			},
		)
		if err != nil {
			return err
		}
		//
		// Create a cluster
		cluster, err := rds.NewCluster(ctx, "aws-classic-postgres-db-cluster",
			&rds.ClusterArgs{
				ClusterIdentifier: pulumi.String("postgres-db-cluster"),
				//AvailabilityZones:        pulumi.ToStringArray(availabilityZones),
				BackupRetentionPeriod:    pulumi.Int(5),
				AllowMajorVersionUpgrade: pulumi.BoolPtr(true),
				DatabaseName:             pulumi.String("example"),
				Engine:                   pulumi.String("aurora-postgresql"),
				EngineVersion:            pulumi.String("12.7"), // aws rds describe-db-engine-versions --engine aurora-postgresql --query '*[].[EngineVersion]' --output text --region aws-region
				MasterPassword:           pulumi.String("dbPassword"),
				MasterUsername:           pulumi.String("mainuser"),
				PreferredBackupWindow:    pulumi.String("02:00-04:00"),
				DbSubnetGroupName:        subnetGroup.Name,
				DeletionProtection:       pulumi.BoolPtr(false),
				SkipFinalSnapshot:        pulumi.BoolPtr(true),
				FinalSnapshotIdentifier:  pulumi.String("final"),
			},
			pulumi.DependsOn([]pulumi.Resource{subnetGroup}),
		)
		if err != nil {
			return err
		}
		/////////////////////////
		// Add instances to support the cluster
		for _, az := range availabilityZones {
			clusterInstName := fmt.Sprintf("aws-classic-cluster-ins-%s", az)
			if _, err := rds.NewClusterInstance(ctx, clusterInstName,
				&rds.ClusterInstanceArgs{
					Identifier:              pulumi.String(clusterInstName),
					ClusterIdentifier:       cluster.ID(),
					InstanceClass:           pulumi.String("db.t3.medium"),
					Engine:                  cluster.Engine,
					EngineVersion:           cluster.EngineVersion,
					AvailabilityZone:        pulumi.StringPtr(az),
					DbParameterGroupName:    pgInst.Name,
					AutoMinorVersionUpgrade: pulumi.BoolPtr(true),
				},
				pulumi.DependsOn([]pulumi.Resource{cluster, pgInst}),
			); err != nil {
				return err
			}
		}
		return nil
	})
}
