function renderResults(publications) {
    if (!publications || publications.length === 0) {
        resultsDiv.innerHTML = '<p>Ничего не найдено.</p>';
        return;
    }
    
    let html = `<h3>Результаты поиска (${publications.length})</h3>`;

    //итерируемся по публикациям
    publications.forEach(pub => {
        const authors = pub.authors ? pub.authors.join(', ') : ''; //превращаем массив в одну строку, разделяя запятой авторов
        const isbn = pub.isbn || '';
        let bbksString = '';
        if (pub.bbks) {
            if (Array.isArray(pub.bbks)) {
                bbksString = pub.bbks.join(' + ');
            } else {
                bbksString = pub.bbks;
            }
        }
        html += `
            <div class="publication-item" data-pub-id="${pub.id}">
                <div class="book-title" style="cursor:pointer; color:#2c3e50;">
                    📘 ${escapeHtml(pub.title)} (${pub.publicationyear})
                </div>
                <div>Авторы: ${escapeHtml(authors)}</div>
                <div>ISBN: ${escapeHtml(isbn)}</div>
                ${bbksString ? `<div>ББК: ${escapeHtml(bbksString)}</div>` : ''}
                <button class="show-libraries-btn" data-pub-id="${pub.id}">📚 Показать библиотеки</button>
                <div class="buildings-list" style="display:none; margin-top:10px; margin-left:20px;"></div>
            </div>
            <hr>
        `;
    });
    resultsDiv.innerHTML = html;

    document.querySelectorAll('.show-libraries-btn').forEach(btn => {
    //итерация по всем кнопкам с классом show-libraries-btn 

        btn.addEventListener('click', (event) => {
            event.preventDefault();
            const pubItem = btn.closest('.publication-item');
            const buildingsDiv = pubItem.querySelector('.buildings-list');
            const pubId = btn.getAttribute('data-pub-id');
            const publication = publications.find(p => p.id == pubId); //находим в массиве publications объект издания, у которого id совпадает с pubId
        
            if (buildingsDiv.style.display === 'none') {
                //если нажали на кнопку показать библиотекеи
                renderBuildings(publication, buildingsDiv);
                buildingsDiv.style.display = 'block';
                btn.textContent = '📖 Скрыть библиотеки';
            } else {
                buildingsDiv.style.display = 'none';
                btn.textContent = '📚 Показать библиотеки';
            }
        });
    });
}

function renderBuildings(publication, container) {
    if (!publication.buildings || publication.buildings.length === 0) {
        container.innerHTML = '<p>Нет информации о наличии книг в библиотеках.</p>';
        return;
    }
    
    let html = '<ul style="list-style-type: none; padding-left: 0;">';
    publication.buildings.forEach(bld => {
        //для каждого здания, в котором есть экземпляры publication выводи инфу
        const avail = bld.availableCopies;
        const total = bld.totalCopies;
        const ratio = `${avail}/${total}`;
        const hasAvailable = avail > 0;

        // Берём первый свободный id экземпляра (хотя лучше брать случайный, но это потом сделаю)
        const copyId = hasAvailable ? bld.availableCopyIds[0] : null;
        
        html += `
            <li style="margin-bottom: 12px; border: 1px solid #eee; padding: 8px;">
                <strong>${escapeHtml(bld.description)}</strong><br>
                Адрес: ${escapeHtml(bld.address)}<br>
                В наличии: ${ratio}<br>
                ${hasAvailable 
                    ? `<button class="reserve-btn" data-copy="${copyId}" data-building-id="${bld.buildingId}">Забронировать</button>` 
                    : '<span style="color: gray;">Нет свободных экземпляров</span>'}
            </li>
        `;
    });
    html += '</ul>';
    container.innerHTML = html;

    container.querySelectorAll('.reserve-btn').forEach(btn => {
        btn.addEventListener('click', async (event) => {
            event.preventDefault();
            const copyId = btn.getAttribute('data-copy');
            if (!copyId) return;
            try {
                const response = await fetch('/reserve', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ copyId: copyId })
                });
                if (!response.ok) {
                    const err = await response.json().catch(() => ({}));
                    throw new Error(err.error || 'Ошибка бронирования');
                }
                const data = await response.json();
                alert(data.message || 'Книга забронирована!');

                const li = btn.closest('li');
                const ratioSpan = li.querySelector('strong ~ br + br + span, strong ~ br + br + div'); // упрощённо

                btn.replaceWith('<span style="color: green;">Забронировано</span>');

            } catch (err) {
                alert(`Ошибка: ${err.message}`);
            }
        });
    });
}