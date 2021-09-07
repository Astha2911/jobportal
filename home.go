package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tutree/jobportal/models"
	"tutree/jobportal/services"
	"tutree/jobportal/utility"

	"github.com/gin-gonic/gin"
)

//GetHome handle home page
func GetHome(c *gin.Context) {
	// role := c.Params.ByName("role")
	role := ""
	tokenString, _ := c.Cookie("tokenString")
	if tokenString != "" {
		role = models.GetRoleByToken(tokenString)
		redirect := "/role/" + role
		c.Redirect(http.StatusFound, redirect)
		return
	}

	// suggestions := models.GetCityCode()
	keywords := models.GetTrendingKeywords(20)
	var trend []string
	for _, key := range keywords {
		trend = append(trend, key.Keyword)
	}
	locString := getLocationFromHeaders(c, "")
	metaTitle := fmt.Sprintf("Search Latest Jobs in USA - Jobs Near Me - %s", utility.GetHost())

	metaDesc := fmt.Sprintf("Jobs Near Me - Find Latest USA Jobs Opportunities on %s. "+
		"Search Delivery Driver, Truck Driver, Caregiver, Part Time or Full Time Jobs. "+
		"Apply Now and Get Hiring Fast.",
		utility.GetHost())

	canonicalLink := utility.GetHostURL()

	metaKeywords := "search jobs, " +
		"fresh jobs, " +
		"essential jobs, " +
		"full time jobs, " +
		"part time jobs, " +
		"hourly jobs, " +
		"jobs in USA"

	ogTags := models.GetHomePageOGTags()

	jobList := GetFilteredHome(locString)
	homeRole := models.GetHomeRole()
	homeIndustry := models.GetHomeIndustry()
	// homeFunction := models.GetHomeFunction()
	homeSkills := models.GetHomeSkills()
	website := services.GetJSONSchemaForWebsite()

	getJobsBySkill := getJobDetailsListBySkill()

	getJobsByOrganization := getJobDetailsListByOrganization()

	c.HTML(http.StatusOK, "home.tmpl.html", addDomainElementsToInterface(gin.H{
		"trending_keywords":        trend,
		"location":                 locString,
		"role":                     role,
		"meta_title":               metaTitle,
		"meta_description":         metaDesc,
		"canonical_link":           canonicalLink,
		"meta_keywords":            metaKeywords,
		"home":                     jobList,
		"home_role":                homeRole,
		"home_industry":            homeIndustry,
		"jobs_by_skills":           getJobsBySkill,
		"jobs_by_organization":     getJobsByOrganization,
		"home_skills":              homeSkills,
		"home_top_skills":          models.GetHomeTopSkills(),
		"home_popular_designation": models.GetPopularDesignationSkills(),
		"ogTags":                   ogTags,
		"website":                  website,
	}))
}

// GetFilteredHome return the strcutred data for Home struct
func GetFilteredHome(locString string) *models.Home {
	// wfh := processJobList(models.GetWorkFromHomeSet(locString))
	org := models.GetCompaniesSet(locString)
	// WfHJob := processJobList(wfh)

	// RJob := processJobList(models.GetRecentJobs(locString))
	// rJobs := processJobList(RJob)

	// WiJ := processJobList(models.GetWalkInJobsSet(locString))
	jobList := &models.Home{
		// WorkHome: &models.WorkFromHome{
		// 	Job: wfh,
		// },
		Companies: org,
		// WorkFromHomeAllLink: "/search/work-from-home-jobs",
		// RecentJob: &models.RecentJobs{
		// 	Job: RJob,
		// },
		// RecentJobLink: "/search/all-jobs",
		// WalkInJobs: &models.WalkInJobs{
		// 	Job: WiJ,
		// },
		// WalkInJobsAllLink:   "/search/walk-in-jobs",
		EmployerChoicesLink: "/search/jobs-by-company",
	}
	return jobList
}

//GetHomeWithRole will handle home page with role
func GetHomeWithRole(c *gin.Context) {
	role := c.Params.ByName("role")

	tokenString, _ := c.Cookie("tokenString")
	locString := getLocationFromHeaders(c, "")
	roleByToken := models.GetRoleByToken(tokenString)

	if role != roleByToken {
		log.Println("Role doesn't match with token string")
		sendErrorPage(http.StatusUnauthorized, "Role does not match", c)
		return
	}

	keywords := models.GetTrendingKeywords(20)
	var trend []string
	for _, key := range keywords {
		trend = append(trend, key.Keyword)
	}
	ogTags := models.GetHomePageOGTags()
	jobList := GetFilteredHome(locString)
	homeSkills := models.GetHomeSkills()
	homeRole := models.GetHomeRole()
	homeIndustry := models.GetHomeIndustry()
	//homeRole := models.GetHomeRole()
	website := services.GetJSONSchemaForWebsite()
	canonicalLink := utility.GetHostURL()

	getJobsBySkill := getJobDetailsListBySkill()

	getJobsByOrganization := getJobDetailsListByOrganization()

	c.HTML(http.StatusOK, "home.tmpl.html", addDomainElementsToInterface(gin.H{
		"trending_keywords":    trend,
		"location":             locString,
		"role":                 role,
		"canonical_link":       canonicalLink,
		"home_industry":        homeIndustry,
		"jobs_by_skills":       getJobsBySkill,
		"jobs_by_organization": getJobsByOrganization,
		// "home_function":         homeFunction,
		"home_skills":              homeSkills,
		"home_top_skills":          models.GetHomeTopSkills(),
		"home_popular_designation": models.GetPopularDesignationSkills(),
		"website":                  website,
		"ogTags":                   ogTags,
		"home_role":                homeRole,
		"home":                     jobList,
	}))
}

//PageNotFound handle 404
func PageNotFound(c *gin.Context) {
	sendErrorPage(http.StatusNotFound, "WE ARE SORRY, BUT THE PAGE YOU REQUESTED WAS NOT FOUND", c)
}

//GetTNC get tnc page
func GetTNC(c *gin.Context) {
	// isAMP := c.GetBool("shouldAMP")

	role := ""
	tokenString, err := c.Cookie("tokenString")

	if err != nil {
		if err != http.ErrNoCookie {
			log.Println("GetTNC: failed:", err)
			sendErrorPage(http.StatusNotFound, "WE ARE SORRY, BUT THE PAGE YOU REQUESTED WAS NOT FOUND", c)
		}
	} else {
		role = models.GetRoleByToken(tokenString)
	}

	tmpl := "tnc.tmpl.html"
	ogTags := models.GetTNCPageOGTags()
	// if isAMP {
	// 	tmpl = "tnc_amp.tmpl.html"
	// }
	c.HTML(200, tmpl, addDomainElementsToInterface(gin.H{
		"amp_tnc": utility.GetAMPTNC(),
		"ogTags":  ogTags,
		"role":    role,
	}))
}

//HandleOldRoute handle old routes - reditect to new
func HandleOldRoute(c *gin.Context) {
	sid := c.Params.ByName("sid")
	state := c.Params.ByName("state")
	city := c.Params.ByName("city")
	zid := c.Params.ByName("zid")

	nid := models.GetNewJobID("", sid)

	c.Redirect(http.StatusFound, utility.GetCanonicalLink(nid, "", state, city, zid))
}

// HandleOldRouteFromTutree handles /online-jobs/ urls
func HandleOldRouteFromTutree(c *gin.Context) {
	sid := c.Params.ByName("sid")
	state := c.Params.ByName("state")
	city := c.Params.ByName("city")
	zid := c.Params.ByName("zid")

	nid, err := strconv.Atoi(sid)
	if err != nil {
		log.Println("HandleOldRouteFromTutree: invalid id:" + sid)
	}

	c.Redirect(http.StatusFound, utility.GetCanonicalLink(nid, "", state, city, zid))
}

// HandleJobByLocationID handles /jobs/?id=1
// We need this because of ancient urls on tutree.com
// the idea is to preserve traffic if possible
func HandleJobByLocationID(c *gin.Context) {
	oid := c.Params.ByName("oid")
	zipcodeID, err := strconv.Atoi(oid)
	if err != nil {
		log.Println("HandleJobByLocationID: invalid id:" + oid + " cannot be converted to int.")
		return
	}
	sid := 245 // this is harcoded, as it came from the old system.
	// jobs/open/%s/%s/%s/%d
	zipcodes, err := models.GetZipcodeByID(zipcodeID)
	if err != nil {
		log.Println("HandleJobByLocationID: model problem:" + oid + " cannot be converted to int.")
		return
	}
	if len(zipcodes) == 0 {
		log.Printf("HandleJobByLocationID: No records found for zipcode:%d\n", zipcodeID)
		return
	}

	c.Redirect(http.StatusFound, utility.GetCanonicalLink(sid, "", zipcodes[0].State, zipcodes[0].CityName, strconv.Itoa(zipcodes[0].ID)))

}

// RedirectTo this will redirect to login page again
func RedirectTo(code int, errMsg string, c *gin.Context, redirect string) {
	c.Redirect(http.StatusFound, redirect)
}

func getJobDetailsListBySkill() models.JobListDetailsBySkill {
	getJobsBySkill := models.JobListDetailsBySkill{
		WarehouseAssociate: models.GetJobBySkill("Warehouse Associate"),
		Driver:             models.GetJobBySkill("Driver"),
		Delivery:           models.GetJobBySkill("Delivery"),
		Nursing:            models.GetJobBySkill("Nursing"),
		Shopper:            models.GetJobBySkill("Shopper"),
	}
	return getJobsBySkill
}

func getJobDetailsListByOrganization() models.JobListDetailsByOrganization {
	getJobsByOrganization := models.JobListDetailsByOrganization{
		Amazon: models.GetJobByOrganization("Amazon"),
		//Uber:     models.GetJobByOrganization("Uber"),
		UberEats: models.GetJobByOrganization("Uber Eats"),
		Doordash: models.GetJobByOrganization("Doordash"),
		
	}

	return getJobsByOrganization
}
