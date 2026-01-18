FROM nginx:alpine

# Copy all files to nginx html directory
COPY . /usr/share/nginx/html/
COPY *.yaml /usr/share/nginx/html/

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"] 