import React from 'react';
import PostEditForm from './components/PostForm';
import { useDispatch, useSelector } from 'react-redux';
import { Skeleton } from 'antd';
import { updatePost, getPost } from '../../actions/posts';
import { useHistory } from 'react-router-dom';
import { useParams } from 'react-router-dom';

function EditPost() {
  const history = useHistory();
  const { id } = useParams();

  const dispatch = useDispatch();

  const { post, loading } = useSelector((state) => {
    const post = state.posts.details[id];

    return {
      post,
      loading: state.posts.loading,
    };
  });

  React.useEffect(() => {
    dispatch(getPost(id));
  }, [id]);

  if (loading) return <Skeleton />;

  const onUpdate = (values) => {
    dispatch(updatePost({ ...post, ...values })).then(() => {
      history.push('/posts');
    });
  };
  return <PostEditForm data={post} onCreate={onUpdate} />;
}

export default EditPost;
