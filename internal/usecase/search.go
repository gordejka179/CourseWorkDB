package usecase

import (
	"github.com/gordejka179/CourseWorkDB/internal/models"
)


type ParametersForm struct {
	Authors []models.Author
	ISBN string
	Title string
	BBKs string
	PublicationYear string
	OtherIndex string
	AdditionalSearch bool
}

type SearchStep struct {
    ShouldRun func(ParametersForm) bool
    Execute func(ParametersForm) (map[int]*PublicationResponse, error)
}


type BuildingAvailability struct {
    BuildingId int `json:"buildingId"`
    Address string `json:"address"`
    Description string `json:"description"`
    TotalCopies int `json:"totalCopies"`
    AvailableCopies int `json:"availableCopies"`
    AvailableInventoryNumbers []string `json:"availableInventoryNumbers"`
}

type PublicationResponse struct {
    Id int `json:"id"`
    Title string `json:"title"`
    PublicationYear int `json:"publicationyear"`
    Authors []models.Author `json:"authors"`
    Isbn string `json:"isbn"`
	BBKs []string `json:"bbks"`
    OtherIndexes []string `json:"otherindexes"`
    Buildings map[int]*BuildingAvailability `json:"buildings"`
}


func (s *Service) SearchPublications(form ParametersForm) (map[int]*PublicationResponse, error) {
    steps := []SearchStep{
        {ShouldRun: func(f ParametersForm) bool { return f.ISBN != "" }, Execute: s.searchByISBN},
		/*
        {ShouldRun: func(f ParametersForm) bool { return f.Title != ""}, Execute: s.searchByTitle},
        {ShouldRun: func(f ParametersForm) bool { return len(f.Authors) != 0 }, Execute: s.searchByAuthors},
        {ShouldRun: func(f ParametersForm) bool { return f.PublicationYear != "" }, Execute: s.searchByPublicationYear},
        {ShouldRun: func(f ParametersForm) bool { return len(f.OtherIndexes) != 0 }, Execute: s.searchByOtherIndexes},
		*/
    }

    for _, step := range steps {
        if step.ShouldRun(form) {
            return step.Execute(form)
        }
    }

	return map[int]*PublicationResponse{}, nil
}

func (s *Service) searchByISBN(form ParametersForm)(map[int]*PublicationResponse, error){
	pubs, err := s.repo.GetPublicationsByISBN(form.ISBN)

	if err != nil{
		return map[int]*PublicationResponse{}, err
	}
	
	pubMap := make(map[int]*PublicationResponse)
	for _, p:= range pubs{
		pubMap[p.ID] = &PublicationResponse{
			Id: p.ID,
			Title: p.Title,
			PublicationYear: p.PublicationYear,
    		Authors: p.Authors,
    		Isbn: p.ISBN,
			BBKs: p.BBKs,
    		OtherIndexes: p.OtherIndexes,
    		Buildings: map[int]*BuildingAvailability{},
		}
	}

	//Поскольку isbn это очень сильный параметр (оставляет очень мало вариантов), 
	//то даже если в запросе будут указаны ещё параметры, то они не будут сужать полученные результаты

	//теперь получаем информацию про конкретные экземпляры
	ids := make([]int, 0, len(pubs))
    for _, pub := range pubs {
        ids = append(ids, pub.ID)
    }

	
	copies, err := s.repo.GetCopiesByIDList(ids)
	if err != nil{
		return map[int]*PublicationResponse{}, err
	}


	for _, copy := range copies{
		pubResp, _ := pubMap[copy.PublicationId]

    	if _, exists := pubResp.Buildings[copy.BuildingId]; !exists {
        	pubResp.Buildings[copy.BuildingId] = &BuildingAvailability{
            	BuildingId: copy.BuildingId,
            	Address:    copy.Address,
            	Description: copy.Description,
            	TotalCopies: 0,
            	AvailableCopies: 0,
            	AvailableInventoryNumbers: []string{},
        	}
    	}
		
		pubMap[copy.PublicationId].Buildings[copy.BuildingId].TotalCopies += 1
		if (copy.ReaderId == 0 && copy.LibrarianId == 0){
			building := pubMap[copy.PublicationId].Buildings[copy.BuildingId]
			building.AvailableCopies += 1
			building.AvailableInventoryNumbers = append(building.AvailableInventoryNumbers, copy.InventoryNumber)
		}
	}


	return pubMap, err
}

