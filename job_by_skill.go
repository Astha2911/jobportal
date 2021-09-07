package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"tutree/jobportal/config"
	"tutree/jobportal/utility"
)

//JobDetailsBySkill struct holds the data for job by skills
type JobDetailsBySkill struct {
	Logo                    string
	LogoAsCompanyName       string
	JobsdetailsTitle        string
	JobsdetailsOrganization string
	JobsdetailsExperience   string
	JobsdetailsLocation     string
	URL                     string
}

////JobDetailsByOrganization struct holds the data for job by company
type JobDetailsByOrganization struct {
	Logo                    string
	JobsdetailsTitle        string
	JobsdetailsOrganization string
	JobsdetailsExperience   string
	JobsdetailsLocation     string
	URL                     string
	LogoAsCompanyName       string
}

type JobListDetailsBySkill struct {
	WarehouseAssociate []JobDetailsBySkill
	Driver             []JobDetailsBySkill
	Delivery           []JobDetailsBySkill
	Nursing            []JobDetailsBySkill
	Shopper            []JobDetailsBySkill
}

type JobListDetailsByOrganization struct {
	Amazon   []JobDetailsByOrganization
	Uber     []JobDetailsByOrganization
	UberEats []JobDetailsByOrganization
	Doordash []JobDetailsByOrganization
	Hogan    []JobDetailsByOrganization
	// OrgLink      []JobDetailsByOrganization
	// Organization []JobDetailsByOrganization
}

//GetJobBySkill function to get the details of job based on skills mentioned
func GetJobBySkill(skill string) []JobDetailsBySkill {
	db, err := config.GetDB2()
	if err != nil {
		log.Println("GetJobBySkill: Failed while connecting with the database :", err)
		return []JobDetailsBySkill{}
	}
	defer db.Close()

	jobdetail := []JobDetailsBySkill{}

	query := `SELECT
		job.id,
		job.job_title,
		org.name,
		org.logo, 
		zip.city_name, 
		zip.state,
		zip.zip
	FROM
		jobs.homepage_job_by_skill skill, 
		jobs.job job, 
		jobs.zipcode zip, 
		jobs.organization org
	WHERE
		skill.is_expired = 0 AND
		skill.job_id = job.id  AND 
		skill.zipcode_id = zip.id AND
		job.organization_id = org.id AND 
		skill.skill = $1
	LIMIT 10`

	rows, err := db.Query(query, skill)
	if err != nil {
		log.Println("GetJobBySkill: failed:", err)
		return jobdetail
	}

	defer rows.Close()

	for rows.Next() {
		var jobID sql.NullInt64
		var title sql.NullString
		var organization sql.NullString
		var logo sql.NullString
		var city sql.NullString
		var state sql.NullString
		var zipcode sql.NullString
		var LogoAsName string

		err := rows.Scan(&jobID, &title, &organization, &logo, &city, &state, &zipcode)
		if err != nil {
			log.Println("GetHomePageJobBySkill: failed:", err)
			continue
		}

		jobId := int(utility.SQLNullIntToInt(jobID))
		Title := utility.SQLNullStringToString(title)
		State := utility.SQLNullStringToString(state)
		City := utility.SQLNullStringToString(city)
		zip := utility.SQLNullStringToString(zipcode)
		Organization := utility.SQLNullStringToString(organization)
		Logo := utility.SQLNullStringToString(logo)

		LogoAsName = utility.GetOrganizationLogoAsName(Organization)

		location := fmt.Sprintf("%s (%s)", City, State)

		url := utility.GetCanonicalLink(
			jobId,
			Title,
			State,
			City,
			zip)
		if utility.GetThemePath() == "htmlJobsdive" {
			url = utility.GetCanonicalLinkForJobsDive(jobId, Title, State, City, zip, Organization)
		}
		data := JobDetailsBySkill{
			JobsdetailsTitle:        Title,
			JobsdetailsOrganization: Organization,
			Logo:                    Logo,
			JobsdetailsExperience:   "1-3",
			JobsdetailsLocation:     location,
			URL:                     url,
			LogoAsCompanyName:       LogoAsName,
		}
		jobdetail = append(jobdetail, data)
	}
	return jobdetail
}

//GetJobByOrganization function to get the details of job based on company name mentioned
func GetJobByOrganization(organization string) []JobDetailsByOrganization {
	db, err := config.GetDB2()
	if err != nil {
		log.Println("GetJobByOrganization: Failed while connecting with the database :", err)
		return nil
	}
	defer db.Close()

	jobdetail := []JobDetailsByOrganization{}

	query := `SELECT
		job.id,
		job.job_title,
		org.name, 
		org.logo,
		zip.city_name, 
		zip.state,
		zip.zip
	FROM
		jobs.homepage_job_by_organization organization, 
		jobs.job job, 
		jobs.zipcode zip, 
		jobs.organization org
	WHERE
		organization.is_expired = 0 AND
		organization.job_id = job.id  AND 
		organization.zipcode_id = zip.id AND
		job.organization_id = org.id AND 
		organization.organization = $1
	LIMIT 10`

	rows, err := db.Query(query, organization)
	if err != nil {
		log.Println("GetJobByOrganization: failed:", err)
		return jobdetail
	}
	defer rows.Close()

	for rows.Next() {
		var jobID sql.NullInt64
		var title sql.NullString
		var organization sql.NullString
		var logo sql.NullString
		var city sql.NullString
		var state sql.NullString
		var zipcode sql.NullString

		err := rows.Scan(&jobID, &title, &organization, &logo, &city, &state, &zipcode)
		if err != nil {
			log.Println("GetJobByOrganization: failed:", err)
			continue
		}

		jobId := int(utility.SQLNullIntToInt(jobID))
		Title := utility.SQLNullStringToString(title)
		State := utility.SQLNullStringToString(state)
		City := utility.SQLNullStringToString(city)
		zip := utility.SQLNullStringToString(zipcode)
		Organization := utility.SQLNullStringToString(organization)
		Logo := utility.SQLNullStringToString(logo)

		LogoAsName := utility.GetOrganizationLogoAsName(Organization)

		location := fmt.Sprintf("%s (%s)", City, State)

		url := utility.GetCanonicalLink(jobId, Title, State, City, zip)
		if utility.GetThemePath() == "htmlJobsdive" {
			url = utility.GetCanonicalLinkForJobsDive(jobId, Title, State, City, zip, Organization)
		}
		data := JobDetailsByOrganization{
			Logo:                    Logo,
			JobsdetailsTitle:        Title,
			JobsdetailsOrganization: Organization,
			JobsdetailsExperience:   "1-3",
			JobsdetailsLocation:     location,
			URL:                     url,
			LogoAsCompanyName:       LogoAsName,
		}
		jobdetail = append(jobdetail, data)
	}
	return jobdetail
}

func JobListDetails(organization string) []JobListDetailsByOrganization {
	var Companies []JobListDetailsByOrganization

	db, err := config.GetDB2()
	if err != nil {
		log.Println("JobListDetails: Failed while connecting with the database :", err)
		return nil
	}
	defer db.Close()

	jobdetail := []JobListDetailsByOrganization{}

	query := `select
		              organization
			 from
				    jobs.homepage_job_by_organization
				where is_organization = true`

	rows, err := db.Query(query)

	if err != nil {
		log.Println("JobListDetailsByOrganization: failed:", err)
		return jobdetail
	}
	defer rows.Close()

	for rows.Next() {

		var organization sql.NullString

		err := rows.Scan(&organization)
		if err != nil {
			log.Println("JobListDetails: failed:", err)
			continue
		}

		Organization := utility.SQLNullStringToString(organization)

		cName := strings.Replace(strings.Replace(strings.Replace(strings.ToLower(Organization), ",", "", -1), " ", "-", -1), ".", "", -1)
		//	orgWords := strings.Split(Name, " ")
		// orgLogoAsName = "CO."
		// if len(orgWords) >= 2 {
		// 	orgLogoAsName = fmt.Sprintf("%s%s", orgWords[0][0:1], orgWords[1][0:1])
		// } else {
		// 	orgLogoAsName = fmt.Sprintf("%s", orgWords[0][0:2])
		// }

		Companies = append(Companies, JobListDetailsByOrganization{
			OrgLink:      os.Getenv("HOST_URL") + "/search/company-" + cName + "-jobs",
			Organization: Organization,
		})
	}
	return Companies
}
