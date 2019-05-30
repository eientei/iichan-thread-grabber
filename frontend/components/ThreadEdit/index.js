import {connect} from "react-redux";
import {withNamespaces} from "../../lib/locale";
import * as React from "react";
import {threadActions} from "../../lib/redux/store";
import {CenteringRowContainer} from "../CenteringRowContainer";
import {WideColumnContainer} from "../WideColumnContainer";
import {RowContainer} from "../RowContainer";
import styled from "@emotion/styled";
import {DragSource, DropTarget} from "react-dnd";
import {ErrorText} from "../ErrorText";
import Router from 'next/router';
import {ColumnContainer} from "../ColumnContainer";

const initialState = {
    stage: 0,
    thread: null,
    cancellable: false,
    cancellation: null,
    queue: [],
    bulk: null,
    setters: {},
};

const HrDot = styled.hr`
  width: 100%;
  border: 0;
  border-bottom: 0.1em #800 dotted;
`;

const Container = styled(RowContainer)`
  flex-grow: 1;
`;

const Left = styled.div`
`;

const Center = styled(RowContainer)`
  flex-grow: 1;
  justify-content: center; 
`;

const Right = styled.div`
  padding-right: 2em;
`;

const SpanCaption = styled.span`
  margin-bottom: 1em;
  display: block;
`;

const Handle = styled.div`
  width: 1em;
`;

const Unsorted = styled.div`
  width: 40%;
  display: flex;
  flex-direction: column;
  align-items: center;
  & > div {
    flex-wrap: wrap;
    justify-content: space-evenly;
    align-items: center;
    flex-grow: 1;
    min-width: 250px;
  }
`;

const Groups = styled.div`
  max-width: 60%;
  flex-grow: 1;
  text-align: center;
  & > div {
    overflow-x: auto;
    overflow-y: hidden;
  }
`;

const Thumbnail = styled.img`
  max-width: 200px;
  max-height: 200px;
`;

const Preview = styled.img`
  max-width: 100%;
  max-height: 100%;
`;

const type = "image";

const BlockA = styled.a`
  display: flex;
  align-items: center;
  justify-content: center;
  ${({expand}) => expand && "width: 200px;"}
  ${({thumb}) => thumb && "height: 200px;"}
`;


const GroupHighlight = styled.div`
  margin: 1em;
  padding: 0.5em;
  min-height: 2em;
  display: flex;
  flex-direction: row;
  align-items: center;
  ${({over}) => over && ("background: #afa;")}
  ${({expand}) => expand && "height: 200px;"}
  ${({border}) => border && "border: 0.2em #800 dashed;"}
`;

const DetailedImageContainer = styled.div`
    display: flex;
    flex-direction: column;
    align-items: center;
`;

const DetailedForm = styled.table`
  width: 50%;
  margin-top: 1em;
  margin-bottom: 1em;
`;

const Th = styled.th`
  text-align: right;
  padding-right: 1em;
`;

const DetailedGroupContainer = styled.div`
  display: flex;
  flex-direction: column;
  margin: 1em;
  padding: 0.5em;
  border: 0.2em #800 dashed;
`;

const InputRadio = styled.input`
  margin: 0 0.3em;
`;

const InputRadioLabel = styled.label`
  font-weight: bold;
`;

const InputText = styled.input`
  width: 100%;
`;

const Md5 = styled.tr`
  height: 4em;
`;

const TagsArea = styled.textarea`
  width: 100%;
  min-height: 8em;
`;

const ImageHighlight = styled.div`
  ${({over}) => over && ("background: #5a5;")}
  ${({expand}) => expand && "min-width: 2em; height: 200px;"}
`;

const ControlContainer = styled(RowContainer)`
  margin-bottom: 1em;
  justify-content: flex-end;
`;

const Image = DragSource(type, {
    canDrag: () => {
        return true;
    },
    beginDrag: ({idx, group}) => ({idx, group}),
}, (connect, monitor) => ({
    connectDragSource: connect.dragSource(),
    isDragging: monitor.isDragging(),
}))(({connectDragSource, file_path, thumb_path, group}) => (
        <BlockA target="_blank" thumb={true} expand={group === 0} href={file_path} ref={instance => connectDragSource(instance)}>
            <Thumbnail src={thumb_path}/>
        </BlockA>
));

const Group = connect()(DropTarget(type, {
    drop: ({dispatch, idx, group}, monitor) => {
        if (monitor.didDrop()) {
            return;
        }
        const item = monitor.getItem();
        dispatch(threadActions.drag(false, {fromgroup: item.group, fromimg: item.idx}, {togroup: idx, toimg: group.length}));
    },
}, (connect, monitor) => ({
    connectDropTarget: connect.dropTarget(),
    isOver: monitor.isOver(),
}))(({connectDropTarget, isOver, group, idx}) => (
        <GroupHighlight over={isOver} border={true} expand={idx > 0} ref={instance => connectDropTarget(instance)}>
            {idx > 0 && (<span>{"(" + group.length + ")"}</span>)}
            <ImageAfter group={idx} idx={0}/>
            {group.map((x,i) => (
                    <React.Fragment key={i}>
                        <Image key={i} group={idx} idx={i} {...x}/>
                        <ImageAfter group={idx} idx={i+1}/>
                    </React.Fragment>
            ))}
        </GroupHighlight>
)));

const GroupAfter = connect()(DropTarget(type, {
    drop: ({dispatch, idx}, monitor) => {
        if (monitor.didDrop()) {
            return;
        }
        const item = monitor.getItem();
        dispatch(threadActions.drag(true, {fromgroup: item.group, fromimg: item.idx}, {togroup: idx, toimg: 0}));
    },
}, (connect, monitor) => ({
    connectDropTarget: connect.dropTarget(),
    isOver: monitor.isOver(),
}))(({connectDropTarget, isOver, idx}) => (
        <GroupHighlight over={isOver} expand={false} ref={instance => connectDropTarget(instance)}>
        </GroupHighlight>
)));

const ImageAfter = connect()(DropTarget(type, {
    drop: ({dispatch, group, idx}, monitor) => {
        if (monitor.didDrop()) {
            return;
        }
        const item = monitor.getItem();
        dispatch(threadActions.drag(false, {fromgroup: item.group, fromimg: item.idx}, {togroup: group, toimg: idx}));
    },
}, (connect, monitor) => ({
    connectDropTarget: connect.dropTarget(),
    isOver: monitor.isOver(),
}))(({connectDropTarget, isOver, group}) => (
        <ImageHighlight over={isOver} expand={group > 0} ref={instance => connectDropTarget(instance)}>
        </ImageHighlight>
)));


const DetailedImage = connect(({thread}) => thread)(withNamespaces(["thread_edit", "error", "common"])(class extends React.Component<*> {
    constructor(props) {
        super(props);
        props.provideSetter(props.group, props.idx, (data) => this.setState(data));
        this.state = ({rating: props.rating, tags: props.tags, parent_md5: props.parent_md5});
    };


    setDeferred(data) {
        const {group, enqueue, idx} = this.props;
        this.setState({...data});
        enqueue({group, idx, ...data});
    }

    setRating(rating) {
        this.setDeferred({rating});
    }

    setParentMd5(parent_md5) {
        this.setDeferred({parent_md5});
    }

    setTags(tags) {
        this.setDeferred({tags});
    }

    render() {
        const {t, file_path, md5} = this.props;
        const {rating, tags, parent_md5} = this.state;
        return (
                <DetailedImageContainer>
                    <BlockA target="_blank" href={file_path}>
                        <Preview src={file_path}/>
                    </BlockA>
                    <DetailedForm>
                        <tbody>
                        <Md5>
                            <Th>{t("md5")}</Th>
                            <td>
                                {md5}
                            </td>
                        </Md5>
                        <tr>
                            <Th>{t("rating")}</Th>
                            <td>
                                <InputRadioLabel onClick={() => this.setRating("e")}>
                                    <InputRadio type="radio" name={md5 + "_rating"} value="e" checked={rating === "e"} onChange={() => null}/>
                                    {t("rating_e")}
                                </InputRadioLabel>
                                <InputRadioLabel onClick={() => this.setRating("q")}>
                                    <InputRadio type="radio" name={md5 + "_rating"} value="q" checked={rating === "q"} onChange={() => null}/>
                                    {t("rating_q")}
                                </InputRadioLabel>
                                <InputRadioLabel onClick={() => this.setRating("s")}>
                                    <InputRadio type="radio" name={md5 + "_rating"} value="s" checked={rating === "s"} onChange={() => null}/>
                                    {t("rating_s")}
                                </InputRadioLabel>
                            </td>
                        </tr>
                        <tr>
                            <Th>
                                <label htmlFor={md5 + "_parent"}>{t("parent_md5")}</label>
                            </Th>
                            <td>
                                <InputText id={md5 + "_parent"} value={parent_md5} onChange={(e) => this.setParentMd5(e.target.value)}/>
                            </td>
                        </tr>
                        <tr>
                            <Th>
                                <label htmlFor={md5 + "_tags"}>{t("tags")}</label>
                            </Th>
                            <td>
                                <TagsArea id={md5 + "_tags"} value={tags} onChange={(e) => this.setTags(e.target.value)}/>
                            </td>
                        </tr>
                        </tbody>
                    </DetailedForm>
                </DetailedImageContainer>
        );
    }
}));


const DetailedGroup = ({t, provideSetter, parentOrder, enqueue, group, idx}) => (
        <DetailedGroupContainer>
            <ControlContainer>
                <button onClick={() => parentOrder(enqueue, group, idx)}>{t("use_parent_order")}</button>
            </ControlContainer>
            {group.map((x,i) => (
                    <React.Fragment key={i}>
                        {i > 0 && (<HrDot/>)}
                        <DetailedImage provideSetter={provideSetter} enqueue={enqueue} group={idx} idx={i} {...x}/>
                    </React.Fragment>
            ))}
        </DetailedGroupContainer>
);

export const ThreadEdit = connect(({thread}) => thread)(withNamespaces(["thread_edit", "error", "common"])(class extends React.Component<*> {
    constructor(props) {
        super(props);
        this.state = initialState;
    };

    provideSetter(group, idx, setter) {
        this.setState(({setters}) => ({setters: {...setters, [group]: {...(setters[group] || {}), [idx]: setter}}}));
    }

    parentOrder(enqueue, group, idx) {
        this.setState(({setters}) => {
            group.forEach((_, i) => {
                if (i > 0) {
                    enqueue({group: idx, idx: i, parent_md5: group[i-1].md5});
                    setters[idx] && setters[idx][i] && setters[idx][i]({parent_md5: group[i-1].md5});
                }
            });
            return null;
        });
    };

    update() {
        const {dispatch} = this.props;
        this.setState(({queue}) => {
            if (queue.length > 0) {
                dispatch(threadActions.update(queue));
            }
            return ({queue: [], bulk: null});
        });
    }

    enqueue(e) {
        setTimeout(() => this.setState(({queue, bulk}) => ({queue: [...queue, e], bulk: bulk == null ? setTimeout(() => this.update(), 10000) : bulk})), 0);
    }

    componentDidMount() {
        const {dispatch} = this.props;
        const path = window.location.pathname.substring(process.env.PUBLIC_PREFIX.length);
        dispatch(threadActions.get(path));
        this.setState({...initialState, thread: path});
    }

    componentWillUnmount() {
        const {dispatch} = this.props;
        const {queue, bulk} = this.state;
        if (bulk != null) {
            clearTimeout(bulk);
        }
        if (queue.length > 0) {
            dispatch(threadActions.update(queue));
        }
    }

    saveThread() {
        const {dispatch} = this.props;
        const {thread} = this.state;
        if (thread === null) {
            return;
        }
        this.setState(({queue}) => {
            if (queue.length > 0) {
                dispatch(threadActions.update(queue));
            }
            dispatch(threadActions.save(thread));
            const cancellation = setTimeout(() => this.setState({cancellable: true}), 5000);
            return ({queue: [], bulk: null, cancellation});
        });
    }

    cancelSave() {
        const {dispatch} = this.props;
        dispatch(threadActions.cancel());
        this.setState({cancellation: null, cancellable: false});
    }

    static getDerivedStateFromProps({submitting}, state) {
        const {cancellation} = state;
        if (!submitting && cancellation != null) {
            clearTimeout(cancellation);
            return ({...state, cancellation: null})
        }
        return null;
    }

    renderStageSort(unsorted, groups) {
        const {t} = this.props;
        return (
                <Container>
                    <Unsorted>
                        <SpanCaption>{t("unsorted_images", {count: unsorted.length})}</SpanCaption>
                        <Group group={unsorted} idx={0}/>
                    </Unsorted>
                    <Handle/>
                    <Groups>
                        <SpanCaption>{t("drag_here", {count: groups.reduce((acc, g) => acc + g.length, 0)})}</SpanCaption>
                        <GroupAfter idx={0}/>
                        {groups.map((g, gi) => (
                                <React.Fragment key={gi}>
                                    <Group group={g} idx={gi+1}/>
                                    <GroupAfter idx={gi+1}/>
                                </React.Fragment>
                        ))}
                    </Groups>
                </Container>
        );
    }

    renderStageFill(groups) {
        const {t} = this.props;
        return (
                <WideColumnContainer>
                    {groups.map((g,i) => (
                            <DetailedGroup t={t} provideSetter={(group, idx, setter) => this.provideSetter(group, idx, setter)} parentOrder={(enqueue, group, idx) => this.parentOrder(enqueue, group, idx)} enqueue={(s) => this.enqueue(s)} key={i+1} idx={i+1} group={g}/>
                    ))}
                </WideColumnContainer>
        );
    }

    renderStage(stage, unsorted, groups) {
        switch(stage) {
            case 0:
                unsorted.sort((a,b) => a.index - b.index);
                return this.renderStageSort(unsorted, groups);
            case 1:
                return this.renderStageFill(groups);
            case 2:
                return null;
            default:
                return "error stage";
        }
    }

    incrementStage() {
        this.setState(({stage, ...state}) => ({...state, stage: stage >= 2 ? stage : stage+1}))
    }

    decrementStage() {
        this.setState(({stage, ...state}) => ({...state, stage: stage > 0 ? stage - 1 : Router.push("/") && 0 || 0}))
    }

    render() {
        const {t, result, submitting, except} = this.props;
        if (result == null) {
            return (
                    <WideColumnContainer>
                        <CenteringRowContainer>
                            {t("fetching")}
                        </CenteringRowContainer>
                    </WideColumnContainer>
            );
        }
        const {stage, cancellable} = this.state;
        const [unsorted, ...groups] = result.groups;
        return (
                <WideColumnContainer>
                    <RowContainer>
                        <Left>
                            {except && <ErrorText>{t("error:" + except)}</ErrorText>}
                        </Left>
                        <Center>
                            <button onClick={() => this.saveThread()} disabled={submitting}>{t("save")}</button>
                            {cancellable && <button onClick={() => this.cancelSave()}>{t("common:cancel")}</button>}
                        </Center>
                        <Right>
                            <button onClick={() => this.decrementStage()}>{t("back")}</button>
                            <button onClick={() => this.incrementStage()}>{t("next")}</button>
                        </Right>
                    </RowContainer>
                    {this.renderStage(stage, unsorted, groups)}
                </WideColumnContainer>
        );
    }
}));