package main

var Links = []struct {
	Href string
	Text string
}{
	{
		Href: "/about",
		Text: "About",
	},
	{
		Href: "/programme-overview",
		Text: "Programme Overview",
	},
	{
		Href: "/keynote-speakers",
		Text: "Keynote Speakers",
	},
	{
		Href: "/registration-and-submission",
		Text: "Registration and submission",
	},
	{
		Href: "/requirements",
		Text: "Requirements",
	},
	{
		Href: "/general-information",
		Text: "General information",
	},
}

type Content struct {
	Title     string
	Paragraph string
}

var homeContent = []Content{
	{
		Title:     "Welcome to AMTC 2022",
		Paragraph: "Admiral Makarov State University of Maritime and Inland Shipping welcomes you to Russia for the International Conference «Arctic: Marine Transportation Challenges – 2022» (AMTC-2022). It is the first time when the problems of the maritime transport and ecology in the Arctic will be discussed globally at the University campus. The Conference is going to become a “smart” platform for discussion and exchange of ideas between Russian and international organizations, executive bodies, business community, research and educational organizations. Apart from main sections, in case of the face-to-face Conference, we will be happy to organize an extensive cultural program in Saint-Petersburg with excursions and informal cocktail. We are looking forward to seeing you as a participant!",
	},
	{
		Title:     "About AMTC-2021",
		Paragraph: "On October 15, 2021, The International Conference «Arctic: Marine Transportation Challenges – 2022» was held at the Admiral Makarov State University of Maritime and Inland Shipping. Due to the restrictions associated with COVID-19 pandemic, the conference sessions were held in mixed offline and online formats. It was held with the support of the Russian Maritime Register of Shipping, Sovcomflot, Administration of the Northern Sea Route, the Arctic and Antarctic Research Institute, with the participation of the leaders of the World Maritime University and the International Association of Maritime Universities. The topics of the plenary reports included the issues of shipping, shipbuilding,port activities, training for work in the Arctic. The conference participants represented Russia, Norway, Finland, USA, Sweden, Japan and China.",
	},
	{
		Title:     "About The Admiral Makarov State University of Maritime and Inland Shipping",
		Paragraph: "The Admiral Makarov State University of Maritime and Inland Shipping is an industry-based vertically-integrated scientific and educational complex. The University incorporates 4 Institutes, College and 8 branches located in different regions of Russia. The main university is situated in Saint-Petersburg and includes 7 campuses. The development of the Northern Sea Route and polar navigational safety is one of the key research activities at the University. The scientific school “Hydrographical support of the Northern Sea Route” is recognized as the leading one. Research and analysis of the Northern Sea Route maritime transport system parameters are being carried out on a regular basis. One example is quantitative and qualitative analysis of a through voyage practice. A separate problem to be solved is the speed mode analysis of high ice-class Arc7 large capacity vessels with deep draught. New researches are devoted to the study of changes in vessels speed modes in ice conditions. For this purpose the Earth remote sensing data, space footage, own calculation methods developed by the University specialists are being used.",
	},
	{
		Title: "Schedule",
		Paragraph: `The conference schedule is as follows:
Full paper for review due
September 15, 2022

Paper acceptance
September 22, 2022

Final Paper Submission
September 26, 2022

Registration
September 26, 2022

Final paper acceptance
September 30, 2022

Conference
October 7-8, 2022`,
	},
}
