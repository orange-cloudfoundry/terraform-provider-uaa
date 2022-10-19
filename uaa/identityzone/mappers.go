package identityzone

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/identityzone/clientsecretpolicyfields"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/identityzone/corsconfigfields"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/identityzone/corsconfignames"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/identityzone/fields"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/identityzone/inputpromptfields"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/identityzone/samlconfigfields"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/identityzone/samlkeyfields"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/identityzone/tokenpolicyfields"
)

func MapIdentityZone(identityZone *api.IdentityZone, data *schema.ResourceData) {

	data.SetId(identityZone.Id)
	data.Set(fields.AccountChooserEnabled.String(), identityZone.Config.AccountChooserEnabled)
	data.Set(fields.IsActive.String(), identityZone.IsActive)
	data.Set(fields.ClientSecretPolicy.String(), mapIdentityZoneClientSecretPolicy(&identityZone.Config.ClientSecretPolicy))
	data.Set(fields.CorsConfig.String(), mapIdentityZoneCorsPolicy(&identityZone.Config.CorsPolicy))
	data.Set(fields.DefaultUserGroups.String(), &identityZone.Config.UserConfig.DefaultGroups)
	data.Set(fields.HomeRedirectUrl.String(), identityZone.Config.Links.HomeRedirect)
	data.Set(fields.IdpDiscoveryEnabled.String(), &identityZone.Config.IdpDiscoveryEnabled)
	data.Set(fields.InputPrompts.String(), mapIdentityZoneInputPrompts(&identityZone.Config.InputPrompts))
	data.Set(fields.IssuerUrl.String(), &identityZone.Config.IssuerUrl)
	data.Set(fields.LogoutRedirectParam.String(), identityZone.Config.Links.Logout.RedirectParameterName)
	data.Set(fields.LogoutRedirectUrl.String(), identityZone.Config.Links.Logout.RedirectUrl)
	data.Set(fields.LogoutAllowedRedirectUrls.String(), identityZone.Config.Links.Logout.AllowedRedirectUrls)
	data.Set(fields.MfaEnabled.String(), identityZone.Config.MfaConfig.IsEnabled)
	data.Set(fields.MfaIdentityProviders.String(), identityZone.Config.MfaConfig.IdentityProviders)
	data.Set(fields.Name.String(), identityZone.Name)
	data.Set(fields.SubDomain.String(), identityZone.SubDomain)
	data.Set(fields.SelfServeEnabled.String(), identityZone.Config.Links.SelfService.Enabled)
	data.Set(fields.SelfServeSignupUrl.String(), identityZone.Config.Links.SelfService.SignupUrl)
	data.Set(fields.SelfServePasswordResetUrl.String(), identityZone.Config.Links.SelfService.PasswordResetUrl)
	data.Set(fields.SamlConfig.String(), mapIdentityZoneSamlConfig(&identityZone.Config.Saml))
	data.Set(fields.TokenPolicy.String(), mapIdentityZoneTokenPolicy(&identityZone.Config.TokenPolicy))
}

func mapIdentityZoneCorsPolicy(data *api.IdentityZoneCorsPolicy) []map[string]interface{} {
	return []map[string]interface{}{
		mapIdentityZoneCorsConfiguration(corsconfignames.Default, &data.DefaultConfiguration),
		mapIdentityZoneCorsConfiguration(corsconfignames.Xhr, &data.XhrConfiguration),
	}
}

func mapIdentityZoneCorsConfiguration(name corsconfignames.CorsConfigName, data *api.IdentityZoneCorsConfig) map[string]interface{} {
	return map[string]interface{}{
		corsconfigfields.AllowedOrigins.String():        data.AllowedOrigins,
		corsconfigfields.AllowedOriginPatterns.String(): data.AllowedOriginPatterns,
		corsconfigfields.AllowedUris.String():           data.AllowedUris,
		corsconfigfields.AllowedUriPatterns.String():    data.AllowedUriPatterns,
		corsconfigfields.AllowedHeaders.String():        data.AllowedHeaders,
		corsconfigfields.AllowedMethods.String():        data.AllowedMethods,
		corsconfigfields.AllowedCredentials.String():    data.AllowedCredentials,
		corsconfigfields.Name.String():                  name.String(),
		corsconfigfields.MaxAge.String():                data.MaxAge,
	}
}

func mapIdentityZoneSamlConfig(data *api.IdentityZoneSamlConfig) []map[string]interface{} {
	return []map[string]interface{}{{
		samlconfigfields.ActiveKeyId.String():              data.ActiveKeyId,
		samlconfigfields.AssertionTtlSeconds.String():      data.AssertionTtlSeconds,
		samlconfigfields.Certificate.String():              data.Certificate,
		samlconfigfields.DisableInResponseToCheck.String(): data.DisableInResponseToCheck,
		samlconfigfields.EntityId.String():                 data.EntityId,
		samlconfigfields.IsAssertionSigned.String():        data.IsAssertionSigned,
		samlconfigfields.IsRequestSigned.String():          data.IsRequestSigned,
		samlconfigfields.Key.String():                      mapIdentityZoneSamlKeys(&data.Keys),
		samlconfigfields.WantAssertionSigned.String():      data.WantAssertionSigned,
		samlconfigfields.WantAuthRequestSigned.String():    data.WantAuthnRequestSigned,
	}}
}

func mapIdentityZoneSamlKeys(data *map[string]api.IdentityZoneSamlKey) (keys []map[string]interface{}) {

	for name, key := range *data {
		keys = append(keys, map[string]interface{}{
			samlkeyfields.Certificate.String(): key.Certificate,
			samlkeyfields.Name.String():        name,
		})
	}

	return keys
}

func mapIdentityZoneClientSecretPolicy(data *api.IdentityZoneClientSecretPolicy) []map[string]interface{} {
	return []map[string]interface{}{{
		clientsecretpolicyfields.MaxLength.String():         data.MaxLength,
		clientsecretpolicyfields.MinDigits.String():         data.MinDigit,
		clientsecretpolicyfields.MinLength.String():         data.MinLength,
		clientsecretpolicyfields.MinLowerCaseChars.String(): data.MinLowerCaseCharacter,
		clientsecretpolicyfields.MinSpecialChars.String():   data.MinSpecialCharacter,
		clientsecretpolicyfields.MinUpperCaseChars.String(): data.MinUpperCaseCharacter,
	}}
}

func mapIdentityZoneTokenPolicy(data *api.IdentityZoneTokenPolicy) []map[string]interface{} {
	return []map[string]interface{}{{
		tokenpolicyfields.AccessTokenTtl.String():       data.AccessTokenTtl,
		tokenpolicyfields.ActiveKeyId.String():          data.ActiveKeyId,
		tokenpolicyfields.IsJwtRevocable.String():       data.IsJwtRevocable,
		tokenpolicyfields.IsRefreshTokenUnique.String(): data.IsRefreshTokenUnique,
		tokenpolicyfields.RefreshTokenFormat.String():   data.RefreshTokenFormat,
		tokenpolicyfields.RefreshTokenTtl.String():      data.RefreshTokenTtl,
	}}
}

func mapIdentityZoneInputPrompts(data *[]api.InputPrompt) (prompts []map[string]interface{}) {

	for _, prompt := range *data {
		prompts = append(prompts, map[string]interface{}{
			inputpromptfields.Name.String():  prompt.Name,
			inputpromptfields.Type.String():  prompt.Type,
			inputpromptfields.Value.String(): prompt.Value,
		})
	}

	return prompts
}
