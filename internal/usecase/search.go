package usecase

import (
	"fmt"
	"strings"

	"github.com/gordejka179/CourseWorkDB/internal/models"
)


type ParametersForm struct {
	Author models.Author
	ISBN string
	Title string
	BBKs []string
	PublicationYear int
	OtherIndexes string
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
    AvailableCopyIds []int `json:"availableCopyIds"`
}

type PublicationResponse struct {
    Id int `json:"id"`
    Title string `json:"title"`
    PublicationYear int `json:"publicationyear"`
    Authors []models.Author `json:"authors"`
    Isbns []string `json:"isbn"`
	BBKs []string `json:"bbks"`
    OtherIndexes []string `json:"otherindexes"`
    Buildings map[int]*BuildingAvailability `json:"buildings"`
}


func (s *Service) GetCopiesByIDList(pubMap map[int]*PublicationResponse, ids []int) (map[int]*PublicationResponse, error){
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
            	AvailableCopyIds: []int{},
        	}
    	}
		
		pubMap[copy.PublicationId].Buildings[copy.BuildingId].TotalCopies += 1
		if (copy.ReaderId == 0 && copy.LibrarianId == 0){
			building := pubMap[copy.PublicationId].Buildings[copy.BuildingId]
			building.AvailableCopies += 1
			building.AvailableCopyIds = append(building.AvailableCopyIds, copy.CopyId)
		}
	}


	return pubMap, err
}



func (s *Service) SearchPublications(form ParametersForm) (map[int]*PublicationResponse, error) {
    steps := []SearchStep{
        {ShouldRun: func(f ParametersForm) bool { return f.ISBN != "" }, Execute: s.searchByISBN},
        {ShouldRun: func(f ParametersForm) bool { return f.Title != ""}, Execute: s.searchByTitle},
        {ShouldRun: func(f ParametersForm) bool { return f.Author.LastName != "" || f.Author.FirstName != "" || f.Author.Patronymic != "" }, Execute: s.searchByAuthor},
		{ShouldRun: func(f ParametersForm) bool { return f.OtherIndexes != "" }, Execute: s.searchByOtherIndexes},
		{ShouldRun: func(f ParametersForm) bool { return len(f.BBKs) != 0 }, Execute: s.searchByBBKs},
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
    		Isbns: p.ISBNs,
			BBKs: p.BBKs,
    		OtherIndexes: p.OtherIndexes,
    		Buildings: map[int]*BuildingAvailability{},
		}
	}

	//Поскольку isbn это очень сильный параметр (оставляет очень мало вариантов), 
	//то даже если в запросе будут указаны ещё параметры, то они не будут сужать полученные результаты

	//теперь получаем информацию про конкретные экземпляры
	ids := make([]int, 0, len(pubMap))
    for id, _ := range pubMap {
        ids = append(ids, id)
    }

	
	return s.GetCopiesByIDList(pubMap, ids)
}


func (s *Service) searchByTitle(form ParametersForm) (map[int]*PublicationResponse, error) {
	pubs, err := s.repo.GetPublicationsByTitle(form.Title)

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
    		Isbns: p.ISBNs,
			BBKs: p.BBKs,
    		OtherIndexes: p.OtherIndexes,
    		Buildings: map[int]*BuildingAvailability{},
		}
	}

	//оставляем только такие экземпляры, где введённые пользователем фамилия + имя + отчество автора принадлежат строке, которая в базе
    authorFormLower := strings.ToLower(form.Author.LastName) + strings.ToLower(form.Author.FirstName) + strings.ToLower(form.Author.Patronymic)
    for id, pub := range pubMap {
        found := false
        for _, author := range pub.Authors {
			authorLower := strings.ToLower(author.LastName) + strings.ToLower(author.FirstName) + strings.ToLower(author.Patronymic)
            if strings.Contains(authorLower, authorFormLower) {
                found = true
                break
            }
        }
        if !found {
            delete(pubMap, id)
        }
    }

	
	//оставляем только такие экземпляры, где введённый пользователем год издания такой же
	year := form.PublicationYear
	if year != 0 {
    	for id, pub := range pubMap {
        	if pub.PublicationYear != year {
            	delete(pubMap, id)
        	}
    	}
	}


	//обработка прочих индексов
	otherIndex := form.OtherIndexes
	if year != 0 {
    	for id, pub := range pubMap {
			found := false
			for _, item := range pub.OtherIndexes {
        		if item == otherIndex {
					found = true
					continue
        		}
			
    		}
        	if found == false{
            	delete(pubMap, id)
        	}
    	}
	}


	//ббк для такого запроса обычно не используется при поиске


	//теперь получаем информацию про конкретные экземпляры
	ids := make([]int, 0, len(pubMap))
    for id, _ := range pubMap {
        ids = append(ids, id)
    }

	return s.GetCopiesByIDList(pubMap, ids)
}


func (s *Service) searchByAuthor(form ParametersForm) (map[int]*PublicationResponse, error) {
	pubs, err := s.repo.GetPublicationsByAuthor(form.Author)

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
    		Isbns: p.ISBNs,
			BBKs: p.BBKs,
    		OtherIndexes: p.OtherIndexes,
    		Buildings: map[int]*BuildingAvailability{},
		}
	}

	
	//оставляем только такие экземпляры, где введённый пользователем год издания такой же
	year := form.PublicationYear
	if year != 0 {
    	for id, pub := range pubMap {
        	if pub.PublicationYear != year {
            	delete(pubMap, id)
        	}
    	}
	}


	//обработка прочих индексов
	otherIndex := form.OtherIndexes
	if year != 0 {
    	for id, pub := range pubMap {
			found := false
			for _, item := range pub.OtherIndexes {
        		if item == otherIndex {
					found = true
					continue
        		}
			
    		}
        	if found == false{
            	delete(pubMap, id)
        	}
    	}
	}


	//ббк для такого запроса обычно не используется при поиске


	//теперь получаем информацию про конкретные экземпляры
	ids := make([]int, 0, len(pubMap))
    for id, _ := range pubMap {
        ids = append(ids, id)
    }

	return s.GetCopiesByIDList(pubMap, ids)
}


func (s *Service) searchByBBKs(form ParametersForm) (map[int]*PublicationResponse, error) {
    fullCodes, err := s.repo.GetFullCodes(form.BBKs)
	fullCodes = append(form.BBKs, fullCodes...)
	if err != nil{
		return map[int]*PublicationResponse{}, fmt.Errorf("ошибка при получении полных кодов ббк: %w", err)
	}

	var allCodes []string
    // Если нужен дополнительный поиск
    if form.AdditionalSearch{
        allCodes, err = s.repo.GetAdditionalCodes(fullCodes)
		if err != nil{
			return map[int]*PublicationResponse{}, fmt.Errorf("ошибка при получении дополнительных кодов ббк: %w", err)
		}
    }
	allCodes = append(fullCodes, allCodes...)
	pubs, err := s.repo.GetPublicationsByBBK(allCodes)
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
    		Isbns: p.ISBNs,
			BBKs: p.BBKs,
    		OtherIndexes: p.OtherIndexes,
    		Buildings: map[int]*BuildingAvailability{},
		}
	}

	//теперь получаем информацию про конкретные экземпляры
	ids := make([]int, 0, len(pubMap))
    for id, _ := range pubMap {
        ids = append(ids, id)
    }

	return s.GetCopiesByIDList(pubMap, ids)
}


func (s *Service) searchByOtherIndexes(form ParametersForm) (map[int]*PublicationResponse, error) {
	pubs, err := s.repo.GetPublicationsByOtherIndex(form.OtherIndexes)
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
    		Isbns: p.ISBNs,
			BBKs: p.BBKs,
    		OtherIndexes: p.OtherIndexes,
    		Buildings: map[int]*BuildingAvailability{},
		}
	}

	//теперь получаем информацию про конкретные экземпляры
	ids := make([]int, 0, len(pubMap))
    for id, _ := range pubMap {
        ids = append(ids, id)
    }

	return s.GetCopiesByIDList(pubMap, ids)
}

