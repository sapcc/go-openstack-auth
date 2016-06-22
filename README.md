# go-openstack-auth
Go openstack auth v3 for getting tokens and endpoints

Example:

		options := auth.AuthV3Options{
			IdentityEndpoint:  viper.GetString(ENV_VAR_AUTH_URL),
			Region:            viper.GetString(ENV_VAR_REGION),
			Username:          viper.GetString(ENV_VAR_USERNAME),
			UserId:            viper.GetString(ENV_VAR_USER_ID),
			Password:          viper.GetString(ENV_VAR_PASSWORD),
			ProjectName:       viper.GetString(ENV_VAR_PROJECT_NAME),
			ProjectId:         viper.GetString(ENV_VAR_PROJECT_ID),
			UserDomainName:    viper.GetString(ENV_VAR_USER_DOMAIN_NAME),
			UserDomainId:      viper.GetString(ENV_VAR_USER_DOMAIN_ID),
			ProjectDomainName: viper.GetString(ENV_VAR_PROJECT_DOMAIN_NAME),
			ProjectDomainId:   viper.GetString(ENV_VAR_PROJECT_DOMAIN_ID),
		}

		authV3 := auth.New(options)
		token, err := authV3.GetToken()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("°°°TOKEN°°°")
		fmt.Println(token.ID)
		fmt.Println("°°°")

		endpoint, err := authV3.GetServiceEndpoint("arc", "staging", "public")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("°°°ENDPOINT°°°")
		fmt.Println(endpoint)
		fmt.Println("°°°")