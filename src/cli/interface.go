package cli

// CVE Generate from https://github.com/CVEProject/cve-schema/releases/tag/v5.1.1 and using https://json2struct.mervine.net/
type CVE struct {
	Containers struct {
		Adp []struct {
			Metrics []struct {
				Other struct {
					Content struct {
						ID      string `json:"id,omitempty"`
						Options []struct {
							Automatable      string `json:"Automatable,omitempty"`
							Exploitation     string `json:"Exploitation,omitempty"`
							Technical_Impact string `json:"Technical Impact,omitempty"`
						} `json:"options,omitempty"`
						Role      string `json:"role,omitempty"`
						Timestamp string `json:"timestamp,omitempty"`
						Version   string `json:"version,omitempty"`
					} `json:"content,omitempty"`
					Type string `json:"type,omitempty"`
				} `json:"other,omitempty"`
			} `json:"metrics,omitempty"`
			ProviderMetadata struct {
				DateUpdated string `json:"dateUpdated,omitempty"`
				OrgID       string `json:"orgId,omitempty"`
				ShortName   string `json:"shortName,omitempty"`
			} `json:"providerMetadata,omitempty"`
			Title string `json:"title,omitempty"`
		} `json:"adp,omitempty"`
		Cna struct {
			Affected []struct {
				Platforms []string `json:"platforms,omitempty"`
				Product   string   `json:"product,omitempty"`
				Vendor    string   `json:"vendor,omitempty"`
				Versions  []struct {
					LessThan    string `json:"lessThan,omitempty"`
					Status      string `json:"status,omitempty"`
					Version     string `json:"version,omitempty"`
					VersionType string `json:"versionType,omitempty"`
				} `json:"versions,omitempty"`
			} `json:"affected,omitempty"`
			CpeApplicability []struct {
				Nodes []struct {
					CpeMatch []struct {
						Criteria              string `json:"criteria,omitempty"`
						VersionEndExcluding   string `json:"versionEndExcluding,omitempty"`
						VersionStartIncluding string `json:"versionStartIncluding,omitempty"`
						Vulnerable            bool   `json:"vulnerable,omitempty"`
					} `json:"cpeMatch,omitempty"`
					Negate   bool   `json:"negate,omitempty"`
					Operator string `json:"operator,omitempty"`
				} `json:"nodes,omitempty"`
			} `json:"cpeApplicability,omitempty"`
			DatePublic   string `json:"datePublic,omitempty"`
			Descriptions []struct {
				Lang  string `json:"lang,omitempty"`
				Value string `json:"value,omitempty"`
			} `json:"descriptions,omitempty"`
			Metrics []struct {
				CvssV3_1 struct {
					BaseScore    float64 `json:"baseScore,omitempty"`
					BaseSeverity string  `json:"baseSeverity,omitempty"`
					VectorString string  `json:"vectorString,omitempty"`
					Version      string  `json:"version,omitempty"`
				} `json:"cvssV3_1,omitempty"`
				Format    string `json:"format,omitempty"`
				Scenarios []struct {
					Lang  string `json:"lang,omitempty"`
					Value string `json:"value,omitempty"`
				} `json:"scenarios,omitempty"`
			} `json:"metrics,omitempty"`
			ProblemTypes []struct {
				Descriptions []struct {
					CweID       string `json:"cweId,omitempty"`
					Description string `json:"description,omitempty"`
					Lang        string `json:"lang,omitempty"`
					Type        string `json:"type,omitempty"`
				} `json:"descriptions,omitempty"`
			} `json:"problemTypes,omitempty"`
			ProviderMetadata struct {
				DateUpdated string `json:"dateUpdated,omitempty"`
				OrgID       string `json:"orgId,omitempty"`
				ShortName   string `json:"shortName,omitempty"`
			} `json:"providerMetadata,omitempty"`
			References []struct {
				Name string   `json:"name,omitempty"`
				Tags []string `json:"tags,omitempty"`
				URL  string   `json:"url,omitempty"`
			} `json:"references,omitempty"`
			Title           string `json:"title,omitempty"`
			XLegacyV4Record struct {
				CVEDataMeta struct {
					Assigner string `json:"ASSIGNER,omitempty"`
					ID       string `json:"ID,omitempty"`
					State    string `json:"STATE,omitempty"`
				} `json:"CVE_data_meta,omitempty"`
				Affects struct {
					Vendor struct {
						VendorData []struct {
							Product struct {
								ProductData []struct {
									ProductName string `json:"product_name,omitempty"`
									Version     struct {
										VersionData []struct {
											VersionValue string `json:"version_value,omitempty"`
										} `json:"version_data,omitempty"`
									} `json:"version,omitempty"`
								} `json:"product_data,omitempty"`
							} `json:"product,omitempty"`
							VendorName string `json:"vendor_name,omitempty"`
						} `json:"vendor_data,omitempty"`
					} `json:"vendor,omitempty"`
				} `json:"affects,omitempty"`
				DataFormat  string `json:"data_format,omitempty"`
				DataType    string `json:"data_type,omitempty"`
				DataVersion string `json:"data_version,omitempty"`
				Description struct {
					DescriptionData []struct {
						Lang  string `json:"lang,omitempty"`
						Value string `json:"value,omitempty"`
					} `json:"description_data,omitempty"`
				} `json:"description,omitempty"`
				Problemtype struct {
					ProblemtypeData []struct {
						Description []struct {
							Lang  string `json:"lang,omitempty"`
							Value string `json:"value,omitempty"`
						} `json:"description,omitempty"`
					} `json:"problemtype_data,omitempty"`
				} `json:"problemtype,omitempty"`
				References struct {
					ReferenceData []struct {
						Name      string `json:"name,omitempty"`
						Refsource string `json:"refsource,omitempty"`
						URL       string `json:"url,omitempty"`
					} `json:"reference_data,omitempty"`
				} `json:"references,omitempty"`
			} `json:"x_legacyV4Record,omitempty"`
		} `json:"cna,omitempty"`
	} `json:"containers,omitempty"`
	CveMetadata struct {
		AssignerOrgID     string `json:"assignerOrgId,omitempty"`
		AssignerShortName string `json:"assignerShortName,omitempty"`
		CveID             string `json:"cveId,omitempty"`
		DatePublished     string `json:"datePublished,omitempty"`
		DateReserved      string `json:"dateReserved,omitempty"`
		DateUpdated       string `json:"dateUpdated,omitempty"`
		State             string `json:"state,omitempty"`
	} `json:"cveMetadata,omitempty"`
	DataType    string `json:"dataType,omitempty"`
	DataVersion string `json:"dataVersion,omitempty"`
}
