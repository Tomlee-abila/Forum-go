//document to handle the user/client like and dislike toggling buttonss
document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('likes-dislikes');

    form.addEventListener('submit', async (event) => {
        event.preventDefault()

        // Get the clicked button and post ID
        const clickedButton = event.submitter
        const postIdInput = form.querySelector('input[name="post_id"]');
        const postId = postIdInput.value

        // Create URL-encoded form data
        const formData = new URLSearchParams()
        formData.append('post_id', postId)
        formData.append('type', clickedButton.value)


        try {
            const response = await fetch('/likes', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: formData.toString()
            })

            if (response.ok) {

                window.location.reload()
            } else {
                // console.error('Server error:', responseText);
                if (response.status === 401) {
                    window.location.href = '/login';
                }
            }
        } catch (error) {
            console.error('Request failed:', error);
        }
    });
});
/*
// Handle like/dislike form submissions
// Like/dislike handler for both posts and comments
document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.like-form').forEach(form => {
        form.addEventListener('submit', async (event) => {
            event.preventDefault();

            const formData = new URLSearchParams();
            formData.append('id', event.target.querySelector('input[name="id"]').value);
            formData.append('item_type', event.target.querySelector('input[name="item_type"]').value);
            formData.append('type', event.submitter.value);

            try {
                const response = await fetch('/likes', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: formData
                });

                if (response.ok) {
                    //error handling for easy debugging of the code
                    const itemType = event.target.querySelector('input[name="item_type"]').value;
                    const itemId = event.target.querySelector('input[name="id"]').value;
                    const countSpan = document.getElementById(`${event.submitter.value}-count-${itemType}-${itemId}`);

                    if (countSpan) {
                        const currentCount = parseInt(countSpan.textContent);
                        // countSpan.textContent = currentCount + 1;
                        window.location.reload()
                    }
                } else if (response.status === 401) {
                    window.location.href = '/login';
                } else {
                    const errorText = await response.text();
                    console.error('Error:', errorText);
                    alert('Failed to process like/dislike');
                }
            } catch (error) {
                console.error('Request failed:', error);
                alert('Failed to process request');
            }
        });
    });
});
*/ 
//if the upper fails use the commented one 
