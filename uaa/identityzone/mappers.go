package identityzone

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/api"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/clientsecretpolicyfields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/configfields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/corsconfigfields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/corsconfignames"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/fields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/samlconfigfields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/samlkeyfields"
	"github.com/terraform-providers/terraform-provider-uaa/uaa/identityzone/tokenpolicyfields"
)

func MapIdentityZone(identityZone *api.IdentityZone, data *schema.ResourceData) {

	data.SetId(identityZone.Id)
	data.Set(fields.Config.String(), mapIdentityZoneConfig(&identityZone.Config))
	data.Set(fields.IsActive.String(), identityZone.IsActive)
	data.Set(fields.Name.String(), identityZone.Name)
	data.Set(fields.SubDomain.String(), identityZone.SubDomain)
	data.Set(fields.ClientSecretPolicy.String(), mapIdentityZoneClientSecretPolicy(&identityZone.Config.ClientSecretPolicy))
}

func mapIdentityZoneConfig(data *api.IdentityZoneConfig) []map[string]interface{} {
	return []map[string]interface{}{{
		configfields.CorsConfig.String():  mapIdentityZoneCorsPolicy(&data.CorsPolicy),
		configfields.Saml.String():        mapIdentityZoneSamlConfig(&data.Saml),
		configfields.TokenPolicy.String(): mapIdentityZoneTokenPolicy(&data.TokenPolicy),
	}}
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
