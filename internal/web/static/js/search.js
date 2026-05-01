function renderResults(publications) {
    if (!publications || publications.length === 0) {
        resultsDiv.innerHTML = '<p>Ничего не найдено.</p>';
        return;
    }
    let html = `<h3>Результаты поиска (${publications.length})</h3>`;
    publications.forEach(pub => {
        const authors = pub.author ? pub.author.join(', ') : '—';
        const isbns = pub.isbn ? pub.isbn.join(', ') : '—';
        html += `
            <div class="book-item" data-pub-id="${pub.id}">
                <div class="book-title">${escapeHtml(pub.title)} (${pub.publicationyear})</div>
                <div>Автор: ${escapeHtml(authors)}</div>
                <div>ISBN: ${escapeHtml(isbns)}</div>
                <div>Доступно в библиотеках:</div>
                <ul>
        `;
        if (pub.buildings && pub.buildings.length > 0) {
            pub.buildings.forEach(bld => {
                const hasCopies = bld.copies && bld.copies.length > 0;
                // кнопка будет, только если есть свободные экземпляры
                const reserveButton = hasCopies
                    ? `<button class="reserve-btn" data-building-id="${bld.id}" data-inventory="${bld.copies[0].inventorynumber}">Забронировать</button>`
                    : '<span style="color: gray;">Нет свободных экземпляров</span>';
                html += `
                    <li>
                        <strong>${escapeHtml(bld.description)}</strong> (${escapeHtml(bld.address || 'адрес не указан')})<br>
                        ${reserveButton}
                    </li>
                `;
            });
        } else {
            html += `<li>Нет информации о наличии книг</li>`;
        }
        html += `</ul></div>`;
    });
    resultsDiv.innerHTML = html;

    // Добавляем обработчики для всех кнопок "Забронировать"
    document.querySelectorAll('.reserve-btn').forEach(btn => {
        btn.addEventListener('click', async (event) => {
            event.preventDefault();
            const inventoryNumber = btn.getAttribute('data-inventory');
            if (!inventoryNumber) return;
            try {
                const response = await fetch('/reserve', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ inventory_number: inventoryNumber })
                });
                if (!response.ok) {
                    const err = await response.json().catch(() => ({}));
                    throw new Error(err.error || 'Ошибка бронирования');
                }
                const data = await response.json();
                alert(data.message || 'Книга забронирована!');
                // можно обновить состояние: убрать кнопку или изменить текст
                btn.replaceWith('<span>Забронировано</span>');
            } catch (err) {
                alert(`Ошибка: ${err.message}`);
            }
        });
    });
}