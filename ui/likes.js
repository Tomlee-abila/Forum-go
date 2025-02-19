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