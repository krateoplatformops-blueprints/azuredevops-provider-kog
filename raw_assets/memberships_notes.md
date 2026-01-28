
# Membership

## Team

team to group

		err = memberships.Create(ctx, e.azCli, memberships.CheckMembershipOptions{
			Organization:        helpers.String(organization),
			SubjectDescriptor:   helpers.String(teamDescriptor),
			ContainerDescriptor: groupDescriptor,
		})

## Group

group to group

		err = memberships.Create(ctx, e.azCli, memberships.CheckMembershipOptions{
			Organization:        helpers.String(organization),
			SubjectDescriptor:   group.Descriptor,
			ContainerDescriptor: containerDescriptor,
		})

## User

user to group
AND
user to team

	for _, descriptor := range groupAndTeamDescriptors {
		err = memberships.Create(ctx, e.azCli, memberships.CheckMembershipOptions{
			Organization:        cr.Spec.Organization,
			SubjectDescriptor:   user.Descriptor,
			ContainerDescriptor: descriptor,
		})
		if err != nil {
			return fmt.Errorf("failed to add user %s to group or team: %w", user.PrincipalName, err)
		}
	}